package main

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/sokoide/ws-ai/pkg/pdf"
)

func drawCenter(pdf *gofpdf.Fpdf, str string) {
	// Get the width of the page
	pageWidth, _ := pdf.GetPageSize()

	// Calculate the width of the text
	textWidth := pdf.GetStringWidth(str)

	// Set X position to center the text
	pdf.SetX((pageWidth - textWidth) / 2)

	// Print the text centered
	pdf.Cell(textWidth, 10, str)
}

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

	pdf.GeneratePdf(title, header, imagePath, user, longText, outputFile)
}
