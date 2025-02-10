package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Define the structure of the playing card

type PlayinCard struct {
	Suit string
	Rank string
}

type Card interface {
	fmt.Stringer
	Name() string
}

type TradingCard struct {
	CollectableName string
}

func NewTradingCard(collectableName string) *TradingCard{
	return &TradingCard{
		CollectableName: collectableName,
	}
}

func (tc *TradingCard) String() string{
	 return tc.CollectableName
}

func(tc *TradingCard) Name() string {
	return tc.String()
}

//func that creates a new card and returns a pointer to it
func NewPlayingCard(suit string, rank string) *PlayinCard{
	playingCard := &PlayinCard{
		Suit: suit,
		Rank: rank,
	}
	return playingCard
}

//add a receiver to the PlayingCard struct
//that is able to stringify the card values
func (pc *PlayinCard) String() string {
	return fmt.Sprintf("%s of %s", pc.Rank, pc.Suit)
}

func(pc *PlayinCard) Name() string {
	return pc.Name()
}


//defination of the deck structure
//reconfigure the deck to support the generic strucure
type Deck [C Card] struct {
	cards []C
}

func NewPlayingCardDeck() *Deck[*PlayinCard] {
	suits := []string {"Diamonds", "Hearts", "Clubs", "Spades"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8","9","10","J", "Q", "K"}

	deck := &Deck[*PlayinCard]{}
	for _, suit := range suits{
		for _, rank := range ranks{
			deck.AddCard(NewPlayingCard(suit, rank))
		}
	}
	
	return deck
}

func (d *Deck [C]) AddCard(card C){
	d.cards = append(d.cards, card)
}

func (d *Deck [C]) RandomCard() C{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	cardIdx := r.Intn(len(d.cards))

	return d.cards[cardIdx]
}

func main(){
	//create a new deck
	deck := NewPlayingCardDeck()

	fmt.Printf("----Drawing a card from the deck------\n")
	card := deck.RandomCard();
	fmt.Printf("------Card drawn of suit:%s-----\n",card )

	//type assertion to check whether the card is of type PlayingCard
	//playingCard, ok := card.(*PlayinCard)

	// if !ok {
	// 	fmt.Printf("card received wasn't a plying card")
	// 	os.Exit(1)
	// }

	fmt.Printf("card suit: %s\n", card.Suit)
	fmt.Printf("card rank: %s\n", card.Rank)
}