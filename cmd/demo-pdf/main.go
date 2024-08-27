package main

import (
	"fmt"

	"github.com/sokoide/ws-ai/pkg/pdf"
)

func main() {
	outputFile := "output.pdf"
	header := "Family Day 2024"
	title := `静かなるホゲの冒険 - A dragon story in town`
	imagePath := "./images/img-sample.png"
	user := `ほげたろう`
	longText := `Here's a dragon story in three chapters, totaling approximately 1000 words:

Chapter 1: The Awakening

In the heart of the Misty Mountains, where ancient peaks pierced the clouds and valleys echoed with secrets, a dragon stirred from its millennial slumber. Azurath, the last of the great sky serpents, slowly opened its eyes, each as large as a shield and glowing with an inner fire. Scales the color of burnished sapphires shifted and gleamed as the dragon stretched its massive body, sending small avalanches of treasure cascading down the sides of its hoard.

For countless centuries, Azurath had slept, guarding its precious collection of gold, gems, and magical artifacts. But now, something had disturbed its rest – a change in the air, a shift in the magical currents that flowed through the earth. The dragon's nostrils flared, taking in scents that had not existed when it first closed its eyes: the acrid tang of smoke from distant human settlements, the metallic odor of forges and industry.

With a rumble that shook the mountain to its foundations, Azurath rose to its feet. Its wings, long folded against its body, unfurled like great sails, stirring up whirlwinds of dust and treasure. The dragon's mind, sharp despite its long sleep, began to process the new world it had awakened to.

As Azurath contemplated its next move, a small figure caught its eye. At the entrance to the vast cavern stood a young woman, her eyes wide with a mixture of fear and wonder. She wore the simple clothes of a village herb-gatherer, a basket of plants clutched tightly to her chest. For a moment, dragon and human locked eyes, mutual curiosity overriding their instinctual reactions.

The girl, Lira, had stumbled upon the dragon's lair while searching for rare medicinal herbs that grew only in the highest reaches of the mountains. She had heard legends of the great dragon Azurath, but like most in her village, she had believed them to be nothing more than old tales told around the hearth on cold winter nights.

Azurath's voice, when it spoke, was like thunder rolling through the cavern. "Child of man, you stand before Azurath, last of the sky serpents. Speak quickly – why have you disturbed my slumber?"

Lira, gathering her courage, stepped forward and bowed deeply. "Great Azurath," she said, her voice trembling but clear, "I meant no disturbance. I am but a simple herb-gatherer, seeking plants to heal my people. I did not know your lair was here."

The dragon's eyes narrowed, assessing the truth in her words. After a moment, it nodded its massive head. "You speak truthfully, little one. But tell me, how long have I slept? What has become of the world beyond these mountains?"

And so, as the first chapter of this new age began, Lira found herself recounting the history of centuries to an ancient dragon, neither of them aware of the profound impact their meeting would have on the world below.

Chapter 2: The Journey

As Lira's tale unfolded, Azurath's expression shifted from curiosity to concern. The world had changed dramatically during its slumber. Kingdoms had risen and fallen, magic had waned, and dragons had passed into legend. Humanity had spread across the land, building cities and harnessing new forms of power.

"I must see this new world for myself," Azurath declared, its voice resonating through the cavern. The dragon turned its gaze to Lira. "You, herb-gatherer, shall be my guide. You will show me what has become of the realms I once knew."

Lira's heart raced at the prospect. Fear and excitement warred within her – to travel with a dragon, to see the world from above! But her village, her family... "Great Azurath," she began hesitantly, "I am honored, but my people will worry if I do not return."

The dragon's laugh was like an avalanche. "Worry not, little one. We shall visit your village first. They will not hinder one under a dragon's protection."`
	longTextJpn := `はい、2000字程度の日本語でドラゴンのお話を作成いたします。

『天龍の守護者』

遠い昔、日本の山奥に「天龍」と呼ばれる神秘的な生き物が住んでいました。天龍は、巨大な翼と鱗に覆われた体を持ち、その姿は威厳に満ちていました。しかし、その存在は長い間、人々の記憶から忘れ去られていました。

山々に囲まれた小さな村、青葉村では、代々「龍守」と呼ばれる家系が天龍の伝説を守り続けていました。現在の龍守は、15歳になったばかりの少女、佐藤美咲でした。美咲は幼い頃から天龍の物語を聞かされて育ちましたが、それが本当に存在するとは信じていませんでした。

ある夏の日、村を襲った大雨により、川が氾濫し、村は大きな被害を受けました。多くの家が流され、畑は水没し、村人たちは絶望的な状況に陥りました。美咲は、何か自分にできることはないかと悩んでいました。

その夜、美咲は不思議な夢を見ました。夢の中で、美咲は山頂に立っていました。そこで、巨大な龍が彼女に語りかけてきたのです。「私は天龍。長い間眠っていたが、今こそ目覚める時が来た。お前は龍守の末裔。私を目覚めさせる力を持っているのだ」

目が覚めた美咲は、夢の意味を理解しようと必死でした。そして、彼女は決心しました。伝説の天龍を探し出し、村を救うために助けを求めることにしたのです。

翌朝早く、美咲は誰にも告げずに家を出ました。彼女は、代々伝わる古い巻物を手に、山へと向かいました。険しい山道を登りながら、美咲は不安と期待が入り混じった気持ちでいっぱいでした。

数時間の登山の末、美咲は巻物に記された神聖な場所に到着しました。そこは、巨大な岩が積み重なってできた洞窟でした。美咲は恐る恐る洞窟に足を踏み入れました。

洞窟の奥に進むにつれ、不思議な光が見えてきました。そして、美咲の目の前に巨大な龍の姿が現れたのです。それは夢で見た天龍そのものでした。

美咲は震える声で天龍に語りかけました。「私は龍守の末裔、佐藤美咲です。私たちの村が大変な危機に陥っています。どうか力を貸してください」

天龍は静かに目を開け、美咲を見つめました。「よく来たな、若き龍守よ。お前の勇気と決意に敬意を表する。しかし、私の力を借りるには代償が必要だ。お前は、その覚悟があるか？」

美咲は一瞬躊躇しましたが、村のことを思い出し、強く頷きました。「はい、どんな代償でも払う覚悟はできています」

天龍は満足げに頷き、こう告げました。「よろしい。その覚悟、確かに受け取った。代償とは、お前がこれからの人生を龍守として生きることだ。天龍と人間の世界の架け橋となり、自然と人間の調和を守り続けることができるか？」

美咲は深く考え、決意を固めました。「はい、私はその役目を全うする覚悟があります」

天龍は大きく翼を広げ、洞窟から飛び立ちました。美咲も龍の背に乗り、村へと向かいました。

村に到着すると、天龍は大きな吐息を吐き出しました。すると、不思議なことに洪水が引き始め、被害を受けた家々や畑が元通りになっていきました。村人たちは驚きと喜びの声を上げ、天を仰ぎ見ました。

その日以来、青葉村は天龍の加護を受け、豊かで平和な暮らしを取り戻しました。美咲は約束通り、龍守としての役目を果たし、村と天龍の絆を深めていきました。

彼女の勇気と決断により、忘れ去られていた伝説が現実となり、人と自然の調和の大切さを人々に教えてくれたのです。

この物語は代々語り継がれ、青葉村の人々の心に刻まれました。そして、美咲の子孫たちもまた、龍守としての使命を受け継いでいくのでした。`

	fmt.Println(len(longText), len(longTextJpn))
	err := pdf.GeneratePdf2(title, header, imagePath, user, longTextJpn, "J1"+outputFile)
	if err != nil {
		fmt.Println(err)
	}

	err = pdf.GeneratePdf2(title+"長い文字ほげホゲホゲ。。もっと長い文字", header, imagePath, user, longTextJpn, "J2"+outputFile)
	if err != nil {
		fmt.Println(err)
	}

	err = pdf.GeneratePdf2(title, header, imagePath, user, longText, "E"+outputFile)
	// err := pdf.GeneratePdf(title, header, imagePath, user, longTextJpn, outputFile)

	if err != nil {
		fmt.Println(err)
	}
}
