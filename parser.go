package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric" // 假设 lyrics 包定义了 Lyrics 结构体

	"github.com/PuerkitoBio/goquery"
	"github.com/go-ego/gse"
	"github.com/samber/lo"
	"golang.org/x/net/html"
)

var seg gse.Segmenter

func init() {
	seg.LoadDict()
}
func ParseTTML(raw string, maindom js.Value) (*lyrics.Lyrics, error) {
	result := lyrics.NewEmptyLyrics()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return nil, err
	}
	meta := doc.Find("tt > metadata").Find("*")
	p := doc.Find("tt > div > p")
	mainAgent, err := parseMetadata(meta, result)
	if err != nil {
		return nil, err
	}
	for _, node := range p.Nodes {
		line, err := parseContent(node, mainAgent, "v1")
		if err != nil {
			return nil, err
		}
		line.Index = len(result.Contents)
		result.Contents = append(result.Contents, line)

	}

	for _, item := range result.Contents {

		item.Ele.Get("style").Set("transition", "none")

		// 以 v 分割Agent
		it := strings.Split(item.Agent, "v")[1]
		// 将item转换为int
		index, err := strconv.Atoi(it)
		if err != nil {
			return nil, errors.New("invalid index")
		}
		if (index % 2) != 0 {
			item.Ele.Get("style").Set("left", "10px")
			item.Ele.Get("style").Set("right", "auto")
			item.Ele.Get("style").Set("transformOrigin", "center left")
		} else {
			item.Ele.Get("classList").Call("add", "right")
			item.Ele.Get("style").Set("left", "auto")
			item.Ele.Get("style").Set("right", "15px")
			item.Ele.Get("style").Set("transformOrigin", "center right")
		}
		item.Primary.Ele.Get("classList").Call("add", "main_lrc")

		segments := seg.Segment([]byte(item.Primary.Ele.Get("innerText").String()))
		// 处理分词结果, 普通模式
		fmt.Println(gse.ToString(segments))
		maindom.Call("append", item.Ele)
		if len(item.Primary.Translates) != 0 {
			item.Ele.Call("append", item.Primary.TranslateEle)
		}
		if len(item.Backgrounds) != 0 {
			bgEle := document.Call("createElement", "div")
			bgEle.Get("classList").Call("add", "bg")
			item.BackgroundsEle = bgEle
			for _, item2 := range item.Backgrounds {
				item2.Ele.Get("classList").Call("add", "text")
				bgEle.Call("append", item2.Ele)
				if len(item2.Translates) != 0 {
					item2.TranslateEle.Get("classList").Call("add", "translation")
					bgEle.Call("append", item2.TranslateEle)
				}

				//if item2.End > item.Primary.End {
				//	item.Primary.End = item2.End
				//}

				item.Ele.Call("append", bgEle)
			}
		}

		for _, item3 := range item.Primary.Blocks {
			if item3.Begin == 0 && item3.End == 0 {
				continue
			}

			backgroundImage, backgroundSize, backgroungPX, _ := generateBackgroundFadeStyle(item3.Ele.Get("offsetWidth").Float(), item3.Ele.Get("offsetHeight").Float(), fadeRatio)
			item3.Ele.Get("style").Set("backgroundImage", backgroundImage)
			item3.Ele.Get("style").Set("backgroundSize", backgroundSize)
			item3.Ele.Get("style").Set("backgroundPositionX", fmt.Sprintf("%vpx", backgroungPX))
			fmt.Println(backgroundImage, backgroundSize, backgroungPX)

			if time.Duration(item3.End-item3.Begin) > time.Duration(1000)*time.Millisecond && len(item3.Text) < 8 {
				text := item3.Ele.Get("innerHTML").String()
				item3.Ele.Set("innerHTML", "")
				for _, ite := range text {
					var blank = document.Call("createElement", "div")
					blank.Set("className", "hl_text")
					blank.Set("innerHTML", js.ValueOf(string(ite)))
					blank.Call("setAttribute", "data-text", string(ite))
					blank.Get("style").Call("setProperty", "--p", "0px")
					item3.Ele.Call("append", blank)
				}
			}
		}

	}
	return result, nil
}

func parseMetadata(meta *goquery.Selection, ly *lyrics.Lyrics) (string, error) {
	agents := lo.Filter(meta.Nodes, func(item *html.Node, index int) bool {
		return item.Type == html.ElementNode && item.Data == "ttm:agent"
	})
	metaNodes := lo.Filter(meta.Nodes, func(item *html.Node, index int) bool {
		return item.Type == html.ElementNode && item.Data == "amll:meta"
	})

	if len(metaNodes) == 0 {
		//return "", errors.New("no metadata found")
		return "v1", nil
	}

	// MainAgent
	mainAgent := "v1"
	mainAgentNode := lo.FindOrElse(agents, nil, func(item *html.Node) bool {
		return getAttr(item, "type") == "person"
	})
	if mainAgentNode != nil {
		mainAgent = getAttrOrElse(mainAgentNode, "xml:id", mainAgent)
	}

	getMetaValues := func(key string) []string {
		return filterEmptyStrings(lo.Map(
			lo.Filter(metaNodes, func(item *html.Node, index int) bool {
				return getAttr(item, "key") == key
			}),
			func(item *html.Node, index int) string {
				return getAttr(item, "value")
			},
		))
	}

	ly.Meta.MusicName = getMetaValues("musicName")
	ly.Meta.Artist = getMetaValues("artists")
	ly.Meta.Album = getMetaValues("album")
	ly.Meta.Isrc = getMetaValues("isrc")

	qmusicId := getMetaValues("qqMusicId")
	if len(qmusicId) > 0 {
		ly.References[lyrics.ReferenceSourceQQMusic] = qmusicId[0]
	}
	spotifyId := getMetaValues("spotifyId")
	if len(spotifyId) > 0 {
		ly.References[lyrics.ReferenceSourceSpotify] = spotifyId[0]
	}
	appleMusicId := getMetaValues("appleMusicId")
	if len(appleMusicId) > 0 {
		ly.References[lyrics.ReferenceSourceAppleMusic] = appleMusicId[0]
	}
	ncmMusicId := getMetaValues("ncmMusicId")
	if len(ncmMusicId) > 0 {
		ly.References[lyrics.ReferenceSourceNetease] = ncmMusicId[0]
	}

	author := getMetaValues("ttmlAuthorGithub")
	authorLogin := getMetaValues("ttmlAuthorGithubLogin")
	if len(author) > 0 || len(authorLogin) > 0 {
		ly.Authors = append(author, authorLogin...)
	}
	return mainAgent, nil
}

func parseContent(p *html.Node, mainAgent, defaultAgent string) (*lyrics.Content, error) {
	begin, err := parseTimestamp(getAttr(p, "begin"))
	if err != nil {
		return nil, fmt.Errorf("invalid begin timestamp: %v", err)
	}
	end, err := parseTimestamp(getAttr(p, "end"))
	if err != nil {
		return nil, fmt.Errorf("invalid end timestamp: %v", err)
	}

	agent := getAttrOrElse(p, "ttm:agent", defaultAgent)
	duet := agent != mainAgent

	lrccitem := document.Call("createElement", "div")
	lrccitem.Set("className", "lyric_item")

	mainlrc := document.Call("createElement", "div")

	lrccitem.Call("append", mainlrc)
	primary := &lyrics.Line{
		Begin:        begin,
		End:          end,
		Translates:   make(map[string]string),
		Duet:         duet,
		Blocks:       make([]*lyrics.Block, 0),
		Ele:          mainlrc,
		TranslateEle: document.Call("createElement", "div"),
	}
	primary.TranslateEle.Get("classList").Call("add", "translation")
	backgrounds := make([]*lyrics.Line, 0)
	for node := range p.ChildNodes() {
		if node.Type == html.ElementNode {
			role := getAttr(node, "ttm:role")
			switch role {
			case "x-translation":
				code := getAttrOrElse(node, "xml:lang", "zh-CN")
				primary.Translates[code] = node.FirstChild.Data
				primary.TranslateEle.Call("append", document.Call("createTextNode", node.FirstChild.Data))

			case "x-roman":
				primary.Translates["romaji"] = node.FirstChild.Data
				//primary.TranslateEle.Call("append", document.Call("createTextNode", node.FirstChild.Data))
			case "x-bg":
				line, err := parseContent(node, mainAgent, agent)
				if err != nil {
					return nil, err
				}
				backgrounds = append(backgrounds, line.Primary)
				break
			default:
				blockBegin, err := parseTimestamp(getAttr(node, "begin"))
				if err != nil {
					return nil, fmt.Errorf("invalid block begin timestamp: %v", err)
				}
				blockEnd, err := parseTimestamp(getAttr(node, "end"))
				if err != nil {
					return nil, fmt.Errorf("invalid block end timestamp: %v", err)
				}

				textele := document.Call("createElement", "div")
				textele.Set("className", "char")
				textele.Get("style").Call("setProperty", "--color", "0.2")
				//textele.Get("style").Call("setProperty", "--p", fmt.Sprintf("%vpx", 0))

				textele.Get("style").Call("setProperty", "--rp", "0%")
				textele.Get("style").Call("setProperty", "--rcolor", "1")
				textele.Get("style").Set("transform", "translateY(10px)")

				text := node.FirstChild.Data

				// 如果第一个字符为(
				if text[0] == '(' {
					// 去掉第一个 (
					text = text[1:]
				}
				if text[len(text)-1] == ')' {
					// 去掉最后一个 )
					text = text[:len(text)-1]
				}

				//if time.Duration(blockEnd-blockBegin) > time.Duration(1200)*time.Millisecond {
				//	for _, ite := range text {
				//		var blank = document.Call("createElement", "div")
				//		blank.Set("className", "hl_text")
				//		blank.Set("innerHTML", js.ValueOf(string(ite)))
				//		textele.Call("append", blank)
				//	}
				//} else {
				textele.Set("innerHTML", text)
				//}

				mainlrc.Call("append", textele)
				primary.Blocks = append(primary.Blocks, &lyrics.Block{
					Begin: blockBegin,
					End:   blockEnd,
					Text:  text,
					Ele:   textele,
				})
			}
		}

		if node.Type == html.TextNode {
			textele := document.Call("createElement", "div")
			textele.Set("className", "char")
			text := node.Data
			// 如果第一个字符为(
			if text[0] == '(' {
				// 去掉第一个 (
				text = text[1:]
			}
			if text[len(text)-1] == ')' {
				// 去掉最后一个 )
				text = text[:len(text)-1]
			}

			if len(node.Data) == 1 {
				textele.Set("innerHTML", "&nbsp;")
			} else {
				textele.Set("innerHTML", text)
			}
			primary.Blocks = append(primary.Blocks, &lyrics.Block{
				Begin: 0,
				End:   0,
				Text:  node.Data,
				Ele:   textele,
			})
			mainlrc.Call("append", textele)
		}
	}
	if primary.Begin == 0 && len(primary.Blocks) > 0 {
		primary.Begin = primary.Blocks[0].Begin
	}
	if primary.End == 0 && len(primary.Blocks) > 0 {
		primary.End = primary.Blocks[len(primary.Blocks)-1].End
	}

	return &lyrics.Content{
		Primary:     primary,
		Backgrounds: backgrounds,
		Ele:         lrccitem,
		Agent:       agent,
	}, nil
}

func getAttr(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func getAttrOrElse(node *html.Node, name, defaultValue string) string {
	value := getAttr(node, name)
	if value == "" {
		return defaultValue
	}
	return value
}

func filterEmptyStrings(list []string) []string {
	return lo.Filter(list, func(item string, index int) bool {
		return item != ""
	})
}

func parseTimestamp(timestamp string) (time.Duration, error) {
	if timestamp == "" {
		return 0, errors.New("timestamp cannot be empty")
	}

	parts := strings.Split(timestamp, ":")
	var hours, minutes int64
	var secondsWithMs string

	switch len(parts) {
	case 3: // HH:MM:SS.MS
		var err error
		hours, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid hours format: %v", err)
		}
		minutes, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes format: %v", err)
		}
		secondsWithMs = parts[2]

	case 2: // MM:SS.MS
		var err error
		minutes, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes format: %v", err)
		}
		secondsWithMs = parts[1]

	case 1: // SS.MS
		secondsWithMs = parts[0]

	default:
		return 0, errors.New("invalid timestamp format")
	}

	secParts := strings.Split(secondsWithMs, ".")
	if len(secParts) == 0 || len(secParts) > 2 {
		return 0, errors.New("invalid timestamp format")
	}

	seconds, err := strconv.ParseInt(secParts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid seconds format: %v", err)
	}

	var milliseconds int64
	if len(secParts) > 1 {
		msString := secParts[1]
		// Handle milliseconds padding
		switch {
		case len(msString) < 3:
			msString = msString + strings.Repeat("0", 3-len(msString))
		case len(msString) > 3:
			msString = msString[:3]
		}

		milliseconds, err = strconv.ParseInt(msString, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid milliseconds format: %v", err)
		}
	}

	// Convert everything to nanoseconds
	totalNanoseconds := (hours*3600+minutes*60+seconds)*1e9 + milliseconds*1e6
	return time.Duration(totalNanoseconds), nil
}

func initLrcBackground(lrc *lyrics.Lyrics) {

	for _, content := range lrc.Contents {
		for _, pitem := range content.Primary.Blocks {
			backgroundImage, backgroundSize, backgroungPX, _ := generateBackgroundFadeStyle(pitem.Ele.Get("offsetWidth").Float(), pitem.Ele.Get("offsetHeight").Float(), fadeRatio)
			pitem.Ele.Get("style").Set("backgroundImage", backgroundImage)
			pitem.Ele.Get("style").Set("backgroundSize", backgroundSize)
			pitem.Ele.Get("style").Set("backgroundPositionX", fmt.Sprintf("%vpx", backgroungPX))
		}
	}
}
