package main

import (
	"log"

	"github.com/meinside/steam-community-market-artifact"
)

const (
	maxNumCardsPerDeck     = 3
	maxNumHeroCardsPerDeck = 1
)

// localized constants
var _localizedHeroes map[artifact.Lang][]string
var _localizedRarities map[artifact.Lang]map[artifact.Rarity]string

func init() {
	// XXX - edit following values when a new hero or localization is added

	_localizedHeroes = map[artifact.Lang][]string{
		artifact.LangEnglish: []string{
			"Axe",
			"Bristleback",
			"Drow Ranger",
			"Kanna",
			"Lich",
			"Tinker",
			"Legion Commander",
			"Lycan",
			"Phantom Assassin",
			"Omniknight",
			"Luna",
			"Bounty Hunter",
			"Ogre Magi",
			"Sniper",
			"Treant Protector",
			"Beastmaster",
			"Enchantress",
			"Sorla Khan",
			"Chen",
			"Zeus",
			"Ursa",
			"Skywrath Mage",
			"Winter Wyvern",
			"Venomancer",
			"Prellex",
			"Earthshaker",
			"Magnus",
			"Sven",
			"Dark Seer",
			"Debbi the Cunning", // basic
			"Mazzie",
			"J'Muy the Wise",       // basic
			"Fahrvhan the Dreamer", // basic
			"Necrophos",
			"Centaur Warrunner",
			"Abaddon",
			"Viper",
			"Timbersaw",
			"Keefe the Bold", // basic
			"Tidehunter",
			"Crystal Maiden",
			"Bloodseeker",
			"Pugna",
			"Lion",
			"Storm Spirit",
			"Meepo",
			"Rix",
			"Outworld Devourer",
			// TODO - add more heroes here
		},
		artifact.LangKorean: []string{
			"도끼전사",
			"가시멧돼지",
			"드로우 레인저",
			"칸나",
			"리치",
			"땜장이",
			"군단 사령관",
			"늑대인간",
			"유령 자객",
			"전능기사",
			"루나",
			"현상금 사냥꾼",
			"오거 마법사",
			"저격수",
			"나무정령 수호자",
			"야수지배자",
			"요술사",
			"솔라 칸",
			"첸",
			"제우스",
			"우르사",
			"하늘분노 마법사",
			"겨울 비룡",
			"맹독사",
			"프렐렉스",
			"지진술사",
			"마그누스",
			"스벤",
			"어둠 현자",
			"교활한 데비", // basic
			"매지",
			"현자 제이무이",              // basic
			"Fahrvhan the Dreamer", // basic
			"강령사제",
			"켄타우로스 전쟁용사",
			"아바돈",
			"바이퍼",
			"벌목꾼",
			"Keefe the Bold", // basic
			"파도사냥꾼",
			"수정의 여인",
			"혈귀",
			"퍼그나",
			"라이온",
			"폭풍령",
			"미포",
			"릭스",
			"외계 침략자",
			// TODO - add more heroes here
		},
		// TODO - add more localizations here
	}

	_localizedRarities = map[artifact.Lang]map[artifact.Rarity]string{
		artifact.LangEnglish: map[artifact.Rarity]string{
			artifact.RarityCommon:   "Common Card",
			artifact.RarityUncommon: "Uncommon Card",
			artifact.RarityRare:     "Rare Card",
		},
		artifact.LangKorean: map[artifact.Rarity]string{
			artifact.RarityCommon:   "일반 카드",
			artifact.RarityUncommon: "고급 카드",
			artifact.RarityRare:     "희귀 카드",
		},
		// TODO - add more localizations here
	}
}

func main() {
	var err error

	language := artifact.LangDefault

	items, err := artifact.FetchAll(artifact.RarityAll, language, artifact.SortColumnName, artifact.SortDirectionAsc)
	if err == nil {
		log.Printf("Number of all items: %d", len(items))

		var numCommons, numCommonCards, priceCommons,
			numUncommons, numUncommonCards, priceUncommons,
			numRares, numRareCards, priceRares int

		for _, item := range items {
			// number of cards per item
			numCards := maxNumCardsPerDeck
			if isHero(item.Name, language) {
				numCards = maxNumHeroCardsPerDeck
			}

			// check rarity
			switch rarityOf(item, language) {
			case artifact.RarityCommon:
				numCommons++
				numCommonCards += numCards
				priceCommons += item.SellPrice * numCards
			case artifact.RarityUncommon:
				numUncommons++
				numUncommonCards += numCards
				priceUncommons += item.SellPrice * numCards
			case artifact.RarityRare:
				numRares++
				numRareCards += numCards
				priceRares += item.SellPrice * numCards
			}
		}

		log.Printf("Price for all %d commons (%d cards): $%.2f", numCommons, numCommonCards, float32(priceCommons)/100.0)
		log.Printf("Price for all %d uncommons (%d cards): $%.2f", numUncommons, numUncommonCards, float32(priceUncommons)/100.0)
		log.Printf("Price for all %d rares (%d cards): $%.2f", numRares, numRareCards, float32(priceRares)/100.0)

		total := float32(priceCommons+priceUncommons+priceRares) / 100.0
		tax := taxOf(total)

		log.Printf("Price for full collection: $%.2f (+ tax $%.2f = $%.2f)", total, tax, total+tax)
	} else {
		log.Printf("* Error while fetching items: %s", err)
	}
}

// check if a card with given name is a hero
func isHero(name string, language artifact.Lang) bool {
	if _, exists := _localizedHeroes[language]; !exists {
		log.Printf("* No heroes defined for language: %s", language)

		return false
	}

	for _, hero := range _localizedHeroes[language] {
		if hero == name {
			return true
		}
	}

	return false
}

// get rarity of given item
func rarityOf(item artifact.MarketItem, language artifact.Lang) artifact.Rarity {
	itemType := item.AssetDescription.Type
	rarities := _localizedRarities[language]

	for k, v := range rarities {
		if itemType == v {
			return k
		}
	}

	return artifact.RarityAll // unknown rarity
}

// calculate tax of given price
func taxOf(price float32) float32 {
	return 0.15 * price
}
