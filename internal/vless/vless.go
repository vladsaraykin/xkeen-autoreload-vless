package vless

import (
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	UUID, Address                                    string
	Port                                             int
	Security, Encryption, Fingerprint, Network, Flow string
	SNI, PublicKey, ShortID, Name                    string
	RawOutbound                                      json.RawMessage
	RawOutbounds                                     []json.RawMessage
	Complex                                          bool
}

type Choice struct {
	Key   string
	Label string
	Special bool
}

func Parse(raw string) (Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return Node{}, err
	}
	if u.Scheme != "vless" {
		return Node{}, fmt.Errorf("unsupported scheme %q", u.Scheme)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return Node{}, fmt.Errorf("invalid port: %w", err)
	}
	q := u.Query()
	n := Node{
		UUID: u.User.Username(), Address: u.Hostname(), Port: port,
		Security:    defaultString(q.Get("security"), "reality"),
		Encryption:  defaultString(q.Get("encryption"), "none"),
		Fingerprint: defaultString(q.Get("fp"), "firefox"),
		Network:     defaultString(q.Get("type"), "tcp"),
		Flow:        defaultString(q.Get("flow"), "xtls-rprx-vision"),
		SNI:         q.Get("sni"), PublicKey: q.Get("pbk"), ShortID: q.Get("sid"),
		Name: u.Fragment,
	}
	if n.UUID == "" || n.Address == "" || n.Port == 0 || n.SNI == "" || n.PublicKey == "" || n.ShortID == "" {
		return Node{}, fmt.Errorf("incomplete VLESS node")
	}
	return n, nil
}

// Select оставлен для обратной совместимости с существующими вызовами.
func Select(nodes []Node, selector, city string) (Node, error) {
	return SelectProfile(nodes, selector, city)
}

// SelectProfile выбирает закрепленный профиль. Если профиль исчез, автоматически
// переключаться на другой сервер нельзя.
func SelectProfile(nodes []Node, country, profile string) (Node, error) {
	for _, n := range nodes {
		if !matchesSelector(n.Name, country) {
			continue
		}
		if profile != "" {
			if n.Name != profile {
				continue
			}
		}
		return n, nil
	}
	return Node{}, fmt.Errorf("профиль %q для группы %q не найден в подписке", profile, country)
}

var countryNames = map[string]string{
	"🇦🇱":"Албания","🇦🇹":"Австрия","🇧🇬":"Болгария","🇧🇷":"Бразилия",
	"🇨🇦":"Канада","🇨🇭":"Швейцария","🇨🇱":"Чили","🇨🇴":"Колумбия",
	"🇨🇾":"Кипр","🇨🇿":"Чехия","🇩🇪":"Германия","🇩🇰":"Дания",
	"🇪🇪":"Эстония","🇪🇸":"Испания","🇫🇮":"Финляндия","🇫🇷":"Франция",
	"🇬🇧":"Великобритания","🇬🇷":"Греция","🇭🇰":"Гонконг","🇭🇷":"Хорватия",
	"🇭🇺":"Венгрия","🇮🇩":"Индонезия","🇮🇪":"Ирландия","🇮🇱":"Израиль",
	"🇮🇳":"Индия","🇮🇹":"Италия","🇯🇵":"Япония","🇰🇿":"Казахстан",
	"🇲🇽":"Мексика","🇲🇾":"Малайзия","🇳🇬":"Нигерия","🇳🇱":"Нидерланды",
	"🇳🇴":"Норвегия","🇳🇿":"Новая Зеландия","🇵🇪":"Перу","🇵🇭":"Филиппины",
	"🇵🇱":"Польша","🇵🇹":"Португалия","🇷🇴":"Румыния","🇷🇸":"Сербия",
	"🇷🇺":"Россия","🇸🇪":"Швеция","🇸🇬":"Сингапур","🇸🇮":"Словения",
	"🇸🇰":"Словакия","🇹🇭":"Таиланд","🇹🇷":"Турция","🇺🇦":"Украина",
	"🇺🇸":"США","🇿🇦":"ЮАР","🇦🇷":"Аргентина","🇦🇺":"Австралия",
}

func Choices(nodes []Node) []Choice {
	grouped := make(map[string]Choice)
	for _, n := range nodes {
		key := selectorKey(n.Name)
		if key == "" { continue }
		if key == "SPECIAL" {
			continue
		}
		if _, ok := grouped[key]; ok { continue }
		label := n.Name
		if known, ok := countryNames[key]; ok {
			label = key + " " + known
		} else if _, country, ok := Location(n.Name); ok {
			label = country
		}
		grouped[key] = Choice{Key:key, Label:label}
	}
	out := make([]Choice, 0, len(grouped)+1)
	for _, c := range grouped { out = append(out,c) }
	sort.Slice(out, func(i,j int) bool {return out[i].Label < out[j].Label})
	for _, n := range nodes {
		if selectorKey(n.Name) == "SPECIAL" {
			out = append(out,Choice{Key:"SPECIAL",Label:"Специальные профили",Special:true})
			break
		}
	}
	return out
}

func Profiles(nodes []Node, selector string) []Node {
	var out []Node
	seen := map[string]bool{}
	for _, n := range nodes {
		if !matchesSelector(n.Name,selector) || seen[n.Name] {continue}
		seen[n.Name]=true
		out=append(out,n)
	}
	sort.Slice(out,func(i,j int)bool{return out[i].Name<out[j].Name})
	return out
}

func selectorKey(name string) string {
	name = strings.TrimSpace(name)
	if flag:=leadingFlag(name); flag!="" {
		if flag=="🇪🇺" {return "SPECIAL"}
		return flag
	}
	if strings.HasPrefix(name,"🤖") {return "SPECIAL"}
	if _,country,ok:=Location(name);ok{return strings.ToLower(country)}
	return strings.ToLower(name)
}

func matchesSelector(name, selector string) bool {
	if selector=="SPECIAL" {return selectorKey(name)=="SPECIAL"}
	if leadingFlag(selector)!="" {return selectorKey(name)==selector}
	return strings.Contains(strings.ToLower(name),strings.ToLower(selector))
}

func leadingFlag(s string) string {
	r:=[]rune(strings.TrimSpace(s))
	if len(r)>=2 && isRegionalIndicator(r[0]) && isRegionalIndicator(r[1]) {return string(r[:2])}
	return ""
}
func isRegionalIndicator(r rune) bool {return r>=0x1F1E6 && r<=0x1F1FF}

func Location(name string) (city, country string, ok bool) {
	parts := strings.Split(name, ",")
	if len(parts) < 2 {
		return "", "", false
	}
	city = strings.TrimSpace(parts[0])
	country = strings.TrimSpace(parts[1])
	if city == "" || country == "" {
		return "", "", false
	}
	return city, country, true
}

func (n Node) StableKey() string {
	if len(n.RawOutbound)>0 {
		var data map[string]any
		if json.Unmarshal(n.RawOutbound,&data)==nil {
			delete(data,"tag")
			if ss,ok:=data["streamSettings"].(map[string]any);ok {
				if ss["security"]=="reality" {
					if rs,ok:=ss["realitySettings"].(map[string]any);ok {
						delete(rs,"serverName")
						delete(rs,"shortId")
					}
				}
			}
			b,_:=json.Marshal(data)
			sum:=sha256.Sum256(b)
			return hex.EncodeToString(sum[:])
		}
	}
	return fmt.Sprintf("%s|%s|%d|%s", n.UUID, n.Address, n.Port, n.PublicKey)
}

func defaultString(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
