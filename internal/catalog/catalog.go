package catalog

type Species struct {
	ID      string
	Label   string
	Profile string
}

type VariantGroup struct {
	ID      string
	LabelEN string
	LabelJA string
}

type Variant struct {
	ID               string
	SpeciesID        string
	BreedOrMorph     string
	Color            string
	PopularityTier   int
	MotionProfile    string
	SourceStatus     string
	LabelEN          string
	LabelJA          string
	SpriteBase       string
	SeedStage        bool
	SourcePath       string
	MotionSourcePath string
	Shape            string
	TintHex          string
	AccentHex        string
}

const (
	SourceStatusDeguMotion     = "degu_motion_source"
	SourceStatusPrototypeOnly  = "prototype_only"
	SourceStatusImageGenQueued = "imagegen_queued"
	SourceStatusMotionDraft    = "motion_source_draft"
	SourceStatusMotionAccepted = "motion_source_accepted"

	MotionProfileDegu               = "degu"
	MotionProfileSmallRodentScurry  = "small-rodent-scurry"
	MotionProfileRabbitHop          = "rabbit-hop"
	MotionProfileDogTrot            = "dog-trot"
	MotionProfileCatStalk           = "cat-stalk"
	MotionProfileGeckoCrawl         = "gecko-crawl"
	MotionProfileTortoisePlod       = "tortoise-plod"
	MotionProfileFerretSlink        = "ferret-slink"
	MotionProfileGuineaPigWaddle    = "guinea-pig-waddle"
	MotionProfileHedgehogShuffle    = "hedgehog-shuffle"
	MotionProfileSquirrelBound      = "squirrel-bound"
	MotionProfileFoxTrot            = "fox-trot"
	MotionProfileRedPandaAmble      = "red-panda-amble"
	MotionProfileOtterSlide         = "otter-slide"
	MotionProfileSugarGliderSkitter = "sugar-glider-skitter"
	MotionProfileCapybaraLumber     = "capybara-lumber"
	MotionProfileSnakeSlither       = "snake-slither"
	MotionProfileDragonPlod         = "dragon-plod"
	MotionProfileFrogHop            = "frog-hop"
	MotionProfileBirdHop            = "bird-hop"
)

var SpeciesList = []Species{
	{ID: "degu", Label: "Degu", Profile: "degu"},
	{ID: "chinchilla", Label: "Chinchilla", Profile: "small-mammal"},
	{ID: "macaroni_mouse", Label: "Macaroni mouse / fat-tailed gerbil", Profile: "small-mammal"},
	{ID: "rabbit", Label: "Rabbit", Profile: "hopper"},
	{ID: "quokka", Label: "Quokka", Profile: "hopper"},
	{ID: "dog", Label: "Small dog", Profile: "companion"},
	{ID: "cat", Label: "Cat", Profile: "companion"},
	{ID: "gecko", Label: "Gecko", Profile: "low-crawler"},
	{ID: "hamster", Label: "Hamster", Profile: "small-mammal"},
	{ID: "ferret", Label: "Ferret", Profile: "small-mammal"},
	{ID: "guinea_pig", Label: "Guinea pig", Profile: "small-mammal"},
	{ID: "hedgehog", Label: "Hedgehog", Profile: "small-mammal"},
	{ID: "squirrel", Label: "Squirrel", Profile: "small-mammal"},
	{ID: "flying_squirrel", Label: "Flying squirrel", Profile: "small-mammal"},
	{ID: "ground_squirrel", Label: "Ground squirrel", Profile: "small-mammal"},
	{ID: "dormouse", Label: "Dormouse", Profile: "small-mammal"},
	{ID: "fox", Label: "Fox", Profile: "companion"},
	{ID: "red_panda", Label: "Red panda", Profile: "companion"},
	{ID: "otter", Label: "Otter", Profile: "low-crawler"},
	{ID: "sugar_glider", Label: "Sugar glider", Profile: "small-mammal"},
	{ID: "capybara", Label: "Capybara", Profile: "large-mammal"},
	{ID: "tortoise", Label: "Tortoise", Profile: "low-crawler"},
	{ID: "rat", Label: "Fancy rat", Profile: "small-mammal"},
	{ID: "mouse", Label: "Fancy mouse", Profile: "small-mammal"},
	{ID: "gerbil", Label: "Mongolian gerbil", Profile: "small-mammal"},
	{ID: "prairie_dog", Label: "Prairie dog", Profile: "small-mammal"},
	{ID: "chipmunk", Label: "Chipmunk", Profile: "small-mammal"},
	{ID: "bearded_dragon", Label: "Bearded dragon", Profile: "low-crawler"},
	{ID: "crested_gecko", Label: "Crested gecko", Profile: "low-crawler"},
	{ID: "corn_snake", Label: "Corn snake", Profile: "low-crawler"},
	{ID: "whites_tree_frog", Label: "White's tree frog", Profile: "hopper"},
	{ID: "salamander", Label: "Salamander", Profile: "low-crawler"},
	{ID: "budgerigar", Label: "Budgerigar", Profile: "bird"},
	{ID: "cockatiel", Label: "Cockatiel", Profile: "bird"},
	{ID: "java_sparrow", Label: "Java sparrow", Profile: "bird"},
	{ID: "parrotlet", Label: "Parrotlet", Profile: "bird"},
	{ID: "lovebird", Label: "Lovebird", Profile: "bird"},
	{ID: "wagtail", Label: "Wagtail", Profile: "bird"},
	{ID: "shoebill", Label: "Shoebill", Profile: "bird"},
}

var VariantGroups = []VariantGroup{
	{ID: "chinchilla", LabelEN: "Chinchillas", LabelJA: "チンチラ"},
	{ID: "hamster", LabelEN: "Hamsters", LabelJA: "ハムスター"},
	{ID: "rabbit", LabelEN: "Rabbits", LabelJA: "うさぎ"},
	{ID: "sugar_glider", LabelEN: "Sugar gliders", LabelJA: "モモンガ"},
	{ID: "small_mammal", LabelEN: "Small animals", LabelJA: "小動物"},
	{ID: "bird", LabelEN: "Birds", LabelJA: "鳥"},
	{ID: "cat", LabelEN: "Cats", LabelJA: "猫"},
	{ID: "dog", LabelEN: "Dogs", LabelJA: "犬"},
	{ID: "reptile_amphibian", LabelEN: "Reptiles and amphibians", LabelJA: "爬虫類・両生類"},
	{ID: "other", LabelEN: "Other animals", LabelJA: "その他"},
}

func VariantGroupByID(id string) (VariantGroup, bool) {
	for _, group := range VariantGroups {
		if group.ID == id {
			return group, true
		}
	}
	return VariantGroup{}, false
}

func VariantGroupIDForSpecies(speciesID string) string {
	switch speciesID {
	case "chinchilla":
		return "chinchilla"
	case "hamster":
		return "hamster"
	case "rabbit":
		return "rabbit"
	case "sugar_glider", "flying_squirrel":
		return "sugar_glider"
	case "budgerigar", "cockatiel", "java_sparrow", "parrotlet", "lovebird", "wagtail", "shoebill":
		return "bird"
	case "cat":
		return "cat"
	case "dog":
		return "dog"
	case "gecko", "bearded_dragon", "crested_gecko", "corn_snake", "whites_tree_frog", "salamander":
		return "reptile_amphibian"
	case "macaroni_mouse", "guinea_pig", "ferret", "hedgehog", "squirrel", "ground_squirrel", "dormouse", "rat", "mouse", "gerbil", "prairie_dog", "chipmunk":
		return "small_mammal"
	default:
		return "other"
	}
}

func VariantGroupForSpecies(speciesID string) VariantGroup {
	if group, ok := VariantGroupByID(VariantGroupIDForSpecies(speciesID)); ok {
		return group
	}
	return VariantGroups[len(VariantGroups)-1]
}

const (
	srcChinchilla                            = "docs/art-source/chinchilla/chinchilla-standard-gray-source-truth-transparent.png"
	srcChinchillaMotion                      = "docs/art-source/chinchilla/motion-source/sheets/chinchilla-standard-gray-source-set00.png"
	srcChinchillaBeige                       = "docs/art-source/chinchilla-beige/motion-source/accepted-frames/set00/frame-00.png"
	srcChinchillaBeigeMotion                 = "docs/art-source/chinchilla-beige/motion-source/sheets/chinchilla-beige-source-set00.png"
	srcChinchillaEbony                       = "docs/art-source/chinchilla-ebony/motion-source/accepted-frames/set00/frame-00.png"
	srcChinchillaEbonyMotion                 = "docs/art-source/chinchilla-ebony/motion-source/sheets/chinchilla-ebony-source-set00.png"
	srcMacaroniMouse                         = "docs/art-source/macaroni-mouse/macaroni-mouse-source-truth.png"
	srcMacaroniMouseMotion                   = "docs/art-source/macaroni-mouse/motion-source/sheets/macaroni-mouse-tan-source-set00.png"
	srcRabbit                                = "docs/source-truth/rabbit-source-truth.png"
	srcRabbitMotion                          = "docs/art-source/rabbit/motion-source/sheets/rabbit-chestnut-agouti-source-set00.png"
	srcQuokka                                = "docs/art-source/quokka/motion-source/accepted-frames/set00/frame-00.png"
	srcQuokkaMotion                          = "docs/art-source/quokka/motion-source/sheets/quokka-source-set00.png"
	srcHimalayanRabbit                       = "docs/art-source/himalayan-rabbit/motion-source/accepted-frames/set00/frame-00.png"
	srcHimalayanRabbitMotion                 = "docs/art-source/himalayan-rabbit/motion-source/sheets/himalayan-rabbit-source-set00.png"
	srcLionheadRabbit                        = "docs/art-source/lionhead-rabbit-brown-white/motion-source/accepted-frames/set00/frame-00.png"
	srcLionheadRabbitMotion                  = "docs/art-source/lionhead-rabbit-brown-white/motion-source/sheets/lionhead-rabbit-brown-white-source-set00.png"
	srcHollandLop                            = "docs/art-source/holland-lop/motion-source/accepted-frames/set00/frame-00.png"
	srcHollandLopMotion                      = "docs/art-source/holland-lop/motion-source/sheets/holland-lop-broken-orange-source-set00.png"
	srcNetherlandDwarf                       = "docs/art-source/netherland-dwarf/motion-source/accepted-frames/set00/frame-00.png"
	srcNetherlandDwarfMotion                 = "docs/art-source/netherland-dwarf/motion-source/sheets/netherland-dwarf-chestnut-source-set00.png"
	srcDog                                   = "docs/art-intake/dog/dog-source-truth-transparent.png"
	srcFrenchBulldog                         = "docs/art-source/french-bulldog/motion-source/accepted-frames/set00/frame-00.png"
	srcFrenchBulldogMotion                   = "docs/art-source/french-bulldog/motion-source/sheets/french-bulldog-fawn-source-set00.png"
	srcToyPoodle                             = "docs/art-source/toy-poodle/motion-source/accepted-frames/set00/frame-00.png"
	srcToyPoodleMotion                       = "docs/art-source/toy-poodle/motion-source/sheets/toy-poodle-apricot-source-set00.png"
	srcMiniatureSchnauzer                    = "docs/art-source/miniature-schnauzer/motion-source/accepted-frames/set00/frame-00.png"
	srcMiniatureSchnauzerMotion              = "docs/art-source/miniature-schnauzer/motion-source/sheets/miniature-schnauzer-salt-pepper-source-set00.png"
	srcCat                                   = "docs/art-source/cat/cat-kijitora-source-truth-transparent.png"
	srcMaineCoon                             = "docs/art-source/maine-coon/motion-source/accepted-frames/set00/frame-00.png"
	srcMaineCoonMotion                       = "docs/art-source/maine-coon/motion-source/sheets/maine-coon-brown-tabby-source-set00.png"
	srcRagdoll                               = "docs/art-source/ragdoll/motion-source/accepted-frames/set00/frame-00.png"
	srcRagdollMotion                         = "docs/art-source/ragdoll/motion-source/sheets/ragdoll-seal-bicolor-source-set00.png"
	srcScottishFold                          = "docs/art-source/scottish-fold/motion-source/accepted-frames/set00/frame-00.png"
	srcScottishFoldMotion                    = "docs/art-source/scottish-fold/motion-source/sheets/scottish-fold-silver-tabby-source-set00.png"
	srcDomesticShorthairCalico               = "docs/art-source/domestic-shorthair/motion-source/accepted-frames/set00/frame-00.png"
	srcDomesticShorthairMotion               = "docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png"
	srcDomesticTabbyWhite                    = "docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/accepted-frames/set00/frame-00.png"
	srcDomesticTabbyWhiteMotion              = "docs/art-source/domestic-shorthair-tabby-white-stocky/motion-source/sheets/domestic-shorthair-tabby-white-stocky-source-set00.png"
	srcBritishShorthair                      = "docs/art-source/british-shorthair/motion-source/accepted-frames/set00/frame-00.png"
	srcBritishShorthairMotion                = "docs/art-source/british-shorthair/motion-source/sheets/british-shorthair-blue-source-set00.png"
	srcMunchkin                              = "docs/art-source/munchkin/motion-source/accepted-frames/set00/frame-00.png"
	srcMunchkinMotion                        = "docs/art-source/munchkin/motion-source/sheets/munchkin-brown-tabby-source-set00.png"
	srcGecko                                 = "docs/art-source/gecko/gecko-source-truth.png"
	srcGeckoMotion                           = "docs/art-source/gecko/motion-source/sheets/gecko-gray-brown-source-set00.png"
	srcLeopardGecko                          = "docs/art-source/leopard-gecko/motion-source/accepted-frames/set00/frame-00.png"
	srcLeopardGeckoMotion                    = "docs/art-source/leopard-gecko/motion-source/sheets/leopard-gecko-source-set00.png"
	srcWhitesTreeFrogBlue                    = "docs/art-source/whites-tree-frog-blue/motion-source/accepted-frames/set00/frame-00.png"
	srcWhitesTreeFrogBlueMotion              = "docs/art-source/whites-tree-frog-blue/motion-source/sheets/whites-tree-frog-blue-source-set00.png"
	srcGuineaPig                             = "docs/art-source/guinea-pig/motion-source/accepted-frames/set00/frame-00.png"
	srcGuineaPigMotion                       = "docs/art-source/guinea-pig/motion-source/sheets/guinea-pig-tricolor-source-set00.png"
	srcGuineaPigRussianSmoke                 = "docs/art-source/guinea-pig-russian-smoke-white/motion-source/accepted-frames/set00/frame-00.png"
	srcGuineaPigRussianMotion                = "docs/art-source/guinea-pig-russian-smoke-white/motion-source/sheets/guinea-pig-russian-smoke-white-source-set00.png"
	srcFancyRat                              = "docs/art-source/fancy-rat/motion-source/accepted-frames/set00/frame-00.png"
	srcFancyRatMotion                        = "docs/art-source/fancy-rat/motion-source/sheets/fancy-rat-hooded-source-set00.png"
	srcChipmunk                              = "docs/art-source/chipmunk/motion-source/accepted-frames/set00/frame-00.png"
	srcChipmunkMotion                        = "docs/art-source/chipmunk/motion-source/sheets/chipmunk-striped-source-set00.png"
	srcAlbinoChipmunk                        = "docs/art-source/albino-chipmunk/motion-source/accepted-frames/set00/frame-00.png"
	srcAlbinoChipmunkMotion                  = "docs/art-source/albino-chipmunk/motion-source/sheets/albino-chipmunk-source-set00.png"
	srcTrueAlbinoChipmunk                    = "docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00/frame-00.png"
	srcTrueAlbinoChipmunkMotion              = "docs/art-source/true-albino-chipmunk/motion-source/sheets/true-albino-chipmunk-source-set00.png"
	srcGroundSquirrel                        = "docs/art-source/richardsons-ground-squirrel/motion-source/accepted-frames/set00/frame-00.png"
	srcGroundSquirrelMotion                  = "docs/art-source/richardsons-ground-squirrel/motion-source/sheets/richardsons-ground-squirrel-source-set00.png"
	srcYorkshireTerrier                      = "docs/art-source/yorkshire-terrier/motion-source/accepted-frames/set00/frame-00.png"
	srcYorkshireTerrierMotion                = "docs/art-source/yorkshire-terrier/motion-source/sheets/yorkshire-terrier-longcoat-source-set00.png"
	srcHamster                               = "docs/art-source/hamster/hamster-source-truth.png"
	srcHamsterMotion                         = "docs/art-source/hamster/motion-source/sheets/hamster-golden-syrian-source-set00.png"
	srcDjungarianHamster                     = "docs/art-source/djungarian-hamster/motion-source/accepted-frames/set00/frame-00.png"
	srcDjungarianHamsterMotion               = "docs/art-source/djungarian-hamster/motion-source/sheets/djungarian-hamster-source-set00.png"
	srcCampbellHamster                       = "docs/art-source/campbell-hamster/motion-source/accepted-frames/set00/frame-00.png"
	srcCampbellHamsterMotion                 = "docs/art-source/campbell-hamster/motion-source/sheets/campbell-hamster-source-set00.png"
	srcRoborovskiHamster                     = "docs/art-source/roborovski-hamster/motion-source/accepted-frames/set00/frame-00.png"
	srcRoborovskiHamsterMotion               = "docs/art-source/roborovski-hamster/motion-source/sheets/roborovski-hamster-source-set00.png"
	srcLonghairHamsterBlackWhite             = "docs/art-source/longhair-hamster-black-white/motion-source/accepted-frames/set00/frame-00.png"
	srcLonghairHamsterBlackWhiteMotion       = "docs/art-source/longhair-hamster-black-white/motion-source/sheets/longhair-hamster-black-white-source-set00.png"
	srcLonghairHamsterBlackWhiteMasked       = "docs/art-source/longhair-hamster-black-white-masked/motion-source/accepted-frames/set00/frame-00.png"
	srcLonghairHamsterBlackWhiteMaskedMotion = "docs/art-source/longhair-hamster-black-white-masked/motion-source/sheets/longhair-hamster-black-white-masked-source-set00.png"
	srcDjungarianHamsterYellow               = "docs/art-source/djungarian-hamster-yellow/motion-source/accepted-frames/set00/frame-00.png"
	srcDjungarianHamsterYellowMotion         = "docs/art-source/djungarian-hamster-yellow/motion-source/sheets/djungarian-hamster-yellow-source-set00.png"
	srcDjungarianHamsterPearlWhite           = "docs/art-source/djungarian-hamster-pearl-white/motion-source/accepted-frames/set00/frame-00.png"
	srcDjungarianHamsterPearlWhiteMotion     = "docs/art-source/djungarian-hamster-pearl-white/motion-source/sheets/djungarian-hamster-pearl-white-source-set00.png"
	srcSugarGlider                           = "docs/art-source/sugar-glider/motion-source/accepted-frames/set00/frame-00.png"
	srcSugarGliderMotion                     = "docs/art-source/sugar-glider/motion-source/sheets/sugar-glider-gray-source-set00.png"
	srcLeucisticSugarGlider                  = "docs/art-source/leucistic-sugar-glider/motion-source/accepted-frames/set00/frame-00.png"
	srcLeucisticSugarGliderMotion            = "docs/art-source/leucistic-sugar-glider/motion-source/sheets/leucistic-sugar-glider-source-set00.png"
	srcAfricanDormouse                       = "docs/art-source/african-dormouse/motion-source/accepted-frames/set00/frame-00.png"
	srcAfricanDormouseMotion                 = "docs/art-source/african-dormouse/motion-source/sheets/african-dormouse-source-set00.png"
	srcAmericanFlyingSquirrel                = "docs/art-source/american-flying-squirrel/motion-source/accepted-frames/set00/frame-00.png"
	srcAmericanFlyingSquirrelMotion          = "docs/art-source/american-flying-squirrel/motion-source/sheets/american-flying-squirrel-source-set00.png"
	srcBudgerigar                            = "docs/art-source/budgerigar/motion-source/accepted-frames/set00/frame-00.png"
	srcBudgerigarMotion                      = "docs/art-source/budgerigar/motion-source/sheets/budgerigar-green-yellow-source-set00.png"
	srcCockatiel                             = "docs/art-source/cockatiel/motion-source/accepted-frames/set00/frame-00.png"
	srcCockatielMotion                       = "docs/art-source/cockatiel/motion-source/sheets/cockatiel-normal-gray-source-set00.png"
	srcJavaSparrow                           = "docs/art-source/java-sparrow/motion-source/accepted-frames/set00/frame-00.png"
	srcJavaSparrowMotion                     = "docs/art-source/java-sparrow/motion-source/sheets/java-sparrow-normal-source-set00.png"
	srcParrotlet                             = "docs/art-source/parrotlet/motion-source/accepted-frames/set00/frame-00.png"
	srcParrotletMotion                       = "docs/art-source/parrotlet/motion-source/sheets/parrotlet-green-source-set00.png"
	srcParrotletBlueGreen                    = "docs/art-source/parrotlet-blue-green/motion-source/accepted-frames/set00/frame-00.png"
	srcParrotletBlueGreenMotion              = "docs/art-source/parrotlet-blue-green/motion-source/sheets/parrotlet-blue-green-source-set00.png"
	srcLovebird                              = "docs/art-source/lovebird/motion-source/accepted-frames/set00/frame-00.png"
	srcLovebirdMotion                        = "docs/art-source/lovebird/motion-source/sheets/lovebird-peach-faced-source-set00.png"
	srcJapaneseGiantSalamander               = "docs/art-source/japanese-giant-salamander/motion-source/accepted-frames/set00/frame-00.png"
	srcJapaneseGiantSalamanderMotion         = "docs/art-source/japanese-giant-salamander/motion-source/sheets/japanese-giant-salamander-source-set00.png"
	srcWhiteWagtail                          = "docs/art-source/white-wagtail/motion-source/accepted-frames/set00/frame-00.png"
	srcWhiteWagtailMotion                    = "docs/art-source/white-wagtail/motion-source/sheets/white-wagtail-source-set00.png"
	srcShoebill                              = "docs/art-source/shoebill-stork/motion-source/accepted-frames/set00/frame-00.png"
	srcShoebillMotion                        = "docs/art-source/shoebill-stork/motion-source/sheets/shoebill-stork-source-set00.png"
	srcNetherlandDwarfHimalayan              = "docs/art-source/netherland-dwarf-himalayan/motion-source/accepted-frames/set00/frame-00.png"
	srcNetherlandDwarfHimalayanMotion        = "docs/art-source/netherland-dwarf-himalayan/motion-source/sheets/netherland-dwarf-himalayan-source-set00.png"
	srcFancyRatBlueHooded                    = "docs/art-source/fancy-rat-blue-hooded/motion-source/accepted-frames/set00/frame-00.png"
	srcFancyRatBlueHoodedMotion              = "docs/art-source/fancy-rat-blue-hooded/motion-source/sheets/fancy-rat-blue-hooded-source-set00.png"
	srcFancyRatChocolateSelf                 = "docs/art-source/fancy-rat-chocolate-self/motion-source/accepted-frames/set00/frame-00.png"
	srcFancyRatChocolateSelfMotion           = "docs/art-source/fancy-rat-chocolate-self/motion-source/sheets/fancy-rat-chocolate-self-source-set00.png"
	srcFancyRatCreamAgouti                   = "docs/art-source/fancy-rat-cream-agouti/motion-source/accepted-frames/set00/frame-00.png"
	srcFancyRatCreamAgoutiMotion             = "docs/art-source/fancy-rat-cream-agouti/motion-source/sheets/fancy-rat-cream-agouti-source-set00.png"
	srcRabbitGray                            = "docs/art-source/rabbit-gray/motion-source/accepted-frames/set00/frame-00.png"
	srcRabbitGrayMotion                      = "docs/art-source/rabbit-gray/motion-source/sheets/rabbit-gray-source-set00.png"
	srcAfricanFatTailedGecko                 = "docs/art-source/african-fat-tailed-gecko/motion-source/accepted-frames/set00/frame-00.png"
	srcAfricanFatTailedGeckoMotion           = "docs/art-source/african-fat-tailed-gecko/motion-source/sheets/african-fat-tailed-gecko-source-set00.png"
	srcFerretSablePanda                      = "docs/art-source/ferret-sable-panda/motion-source/accepted-frames/set00/frame-00.png"
	srcFerretSablePandaMotion                = "docs/art-source/ferret-sable-panda/motion-source/sheets/ferret-sable-panda-source-set00.png"
	srcFerretSable                           = "docs/art-source/ferret-sable/motion-source/accepted-frames/set00/frame-00.png"
	srcFerretSableMotion                     = "docs/art-source/ferret-sable/motion-source/sheets/ferret-sable-source-set00.png"
	srcFerretAlbino                          = "docs/art-source/ferret-albino/motion-source/accepted-frames/set00/frame-00.png"
	srcFerretAlbinoMotion                    = "docs/art-source/ferret-albino/motion-source/sheets/ferret-albino-source-set00.png"
)

var Variants = []Variant{
	deguVariant("wild_agouti", "wild agouti", "degu_wild_agouti"),
	deguVariant("black", "black", "degu_black"),
	deguVariant("blue", "blue (slate gray)", "degu_blue"),
	deguVariant("gray", "gray", "degu_gray"),
	deguVariant("white_cream", "white / cream", "degu_white_cream"),
	deguVariant("sand_champagne", "sand / champagne", "degu_sand_champagne"),
	deguVariant("chocolate", "chocolate", "degu_chocolate"),
	deguVariant("black_pied", "black pied", "degu_black_pied"),
	deguVariant("agouti_pied", "agouti pied", "degu_agouti_pied"),
	deguVariant("blue_pied", "blue pied (slate gray)", "degu_blue_pied"),
	deguVariant("cream_pied", "cream pied", "degu_cream_pied"),

	acceptedMotionVariant("chinchilla_standard_gray", "chinchilla", "Chinchilla", "チンチラ", "chinchilla_standard_gray", srcChinchilla, srcChinchillaMotion, "", "", "Chinchilla", "standard gray", 1),
	acceptedMotionVariant("chinchilla_beige", "chinchilla", "Chinchilla - beige", "チンチラ（ベージュ）", "chinchilla_beige", srcChinchillaBeige, srcChinchillaBeigeMotion, "", "", "Chinchilla", "beige", 2),
	acceptedMotionVariant("chinchilla_ebony", "chinchilla", "Chinchilla - ebony", "チンチラ（エボニー）", "chinchilla_ebony", srcChinchillaEbony, srcChinchillaEbonyMotion, "", "", "Chinchilla", "ebony", 2),
	sourceVariant("chinchilla_white_mosaic", "chinchilla", "Chinchilla - white mosaic", "chinchilla_white_mosaic", srcChinchilla, "e6e1d3", "8f9690"),

	acceptedMotionVariant("macaroni_mouse_tan", "macaroni_mouse", "Macaroni mouse", "マカロニマウス", "macaroni_mouse_tan", srcMacaroniMouse, srcMacaroniMouseMotion, "", "", "Macaroni mouse", "tan", 1),
	sourceVariant("macaroni_mouse_gray", "macaroni_mouse", "Macaroni mouse - gray", "macaroni_mouse_gray", srcMacaroniMouse, "8f8f84", ""),
	sourceVariant("macaroni_mouse_cream", "macaroni_mouse", "Macaroni mouse - cream", "macaroni_mouse_cream", srcMacaroniMouse, "dcc698", ""),

	acceptedMotionVariant("rabbit_chestnut_agouti", "rabbit", "Rabbit", "うさぎ", "rabbit_chestnut_agouti", srcRabbit, srcRabbitMotion, "", "", "Rabbit", "chestnut agouti", 1),
	sourceVariant("rabbit_black", "rabbit", "Rabbit - black", "rabbit_black", srcRabbit, "2f2c29", ""),
	sourceVariant("rabbit_white", "rabbit", "Rabbit - white", "rabbit_white", srcRabbit, "ebe7dc", ""),
	sourceVariant("rabbit_blue_gray", "rabbit", "Rabbit - blue gray", "rabbit_blue_gray", srcRabbit, "777f86", ""),
	sourceVariant("rabbit_fawn", "rabbit", "Rabbit - fawn", "rabbit_fawn", srcRabbit, "bf8754", ""),
	acceptedMotionVariant("quokka", "quokka", "Quokka", "クアッカワラビー", "quokka", srcQuokka, srcQuokkaMotion, "a65f25", "d48a3d", "Quokka", "warm brown", 2),

	sourceVariant("dog_cream_tan", "dog", "Small dog - cream and tan", "dog_cream_tan", srcDog, "", ""),
	sourceVariant("dog_black_tan", "dog", "Small dog - black and tan", "dog_black_tan", srcDog, "3a3026", "b8874c"),
	sourceVariant("dog_white", "dog", "Small dog - white", "dog_white", srcDog, "e8e1d0", ""),
	sourceVariant("dog_sable", "dog", "Small dog - sable", "dog_sable", srcDog, "9a6839", ""),
	sourceVariant("dog_gray", "dog", "Small dog - gray", "dog_gray", srcDog, "8d8d88", ""),

	sourceVariant("cat_brown_tabby", "cat", "Cat - brown tabby", "cat_brown_tabby", srcCat, "", ""),
	sourceVariant("cat_black", "cat", "Cat - black", "cat_black", srcCat, "252321", ""),
	sourceVariant("cat_white", "cat", "Cat - white", "cat_white", srcCat, "eee8dd", ""),
	sourceVariant("cat_orange_tabby", "cat", "Cat - orange tabby", "cat_orange_tabby", srcCat, "c87830", ""),
	sourceVariant("cat_gray", "cat", "Cat - gray", "cat_gray", srcCat, "818783", ""),
	sourceVariant("cat_cream", "cat", "Cat - cream", "cat_cream", srcCat, "d7b983", ""),

	acceptedMotionVariant("gecko_gray_brown", "gecko", "Gecko", "ヤモリ", "gecko_gray_brown", srcGecko, srcGeckoMotion, "", "", "Gecko", "gray brown", 1),
	acceptedMotionVariant("gecko_leopard", "gecko", "Leopard gecko", "ヒョウモントカゲモドキ", "gecko_leopard", srcLeopardGecko, srcLeopardGeckoMotion, "d9a83c", "2f2418", "Leopard gecko", "yellow tan spotted", 2),
	sourceVariant("gecko_tangerine", "gecko", "Gecko - tangerine", "gecko_tangerine", srcGecko, "de8a34", ""),
	sourceVariant("gecko_blizzard", "gecko", "Gecko - blizzard", "gecko_blizzard", srcGecko, "d8d2c5", ""),
	sourceVariant("gecko_albino", "gecko", "Gecko - albino", "gecko_albino", srcGecko, "e0cfa8", ""),

	acceptedMotionVariant("hamster_golden_syrian", "hamster", "Hamster", "ハムスター", "hamster_golden_syrian", srcHamster, srcHamsterMotion, "", "", "Hamster", "golden Syrian", 1),
	acceptedMotionVariant("djungarian_hamster", "hamster", "Djungarian hamster", "ジャンガリアンハムスター", "djungarian_hamster", srcDjungarianHamster, srcDjungarianHamsterMotion, "a9a49a", "f4efe7", "Djungarian hamster", "winter white gray", 2),
	acceptedMotionVariant("campbell_hamster", "hamster", "Campbell hamster", "キャンベルハムスター", "campbell_hamster", srcCampbellHamster, srcCampbellHamsterMotion, "9b846f", "d6c9b8", "Campbell dwarf hamster", "warm gray brown", 2),
	acceptedMotionVariant("roborovski_hamster", "hamster", "Roborovski hamster", "ロボロフスキー", "roborovski_hamster", srcRoborovskiHamster, srcRoborovskiHamsterMotion, "d7ad73", "f3eadc", "Roborovski hamster", "sandy white", 2),
	acceptedMotionVariant("longhair_hamster_black_white", "hamster", "Longhair hamster - black white", "白黒長毛ハムスター", "longhair_hamster_black_white", srcLonghairHamsterBlackWhite, srcLonghairHamsterBlackWhiteMotion, "f0e8db", "2f2c2b", "Longhair hamster", "black white", 2),
	acceptedMotionVariant("longhair_hamster_black_white_masked", "hamster", "Longhair hamster - black white masked", "白黒長毛ハムスター2", "longhair_hamster_black_white_masked", srcLonghairHamsterBlackWhiteMasked, srcLonghairHamsterBlackWhiteMaskedMotion, "f1ede6", "252423", "Longhair hamster", "black white masked", 2),
	acceptedMotionVariant("djungarian_hamster_yellow", "hamster", "Yellow Djungarian hamster", "イエロージャンガリアン", "djungarian_hamster_yellow", srcDjungarianHamsterYellow, srcDjungarianHamsterYellowMotion, "dcc08a", "f0e0bd", "Djungarian hamster", "yellow", 2),
	acceptedMotionVariant("djungarian_hamster_pearl_white", "hamster", "Pearl white Djungarian hamster", "パールホワイトジャンガリアン", "djungarian_hamster_pearl_white", srcDjungarianHamsterPearlWhite, srcDjungarianHamsterPearlWhiteMotion, "eee9df", "a69d92", "Djungarian hamster", "pearl white", 2),
	sourceVariant("hamster_cream", "hamster", "Hamster - cream", "hamster_cream", srcHamster, "e2c18b", ""),
	sourceVariant("hamster_black_banded", "hamster", "Hamster - black banded", "hamster_black_banded", srcHamster, "3f342d", "ddd2bd"),
	sourceVariant("hamster_white", "hamster", "Hamster - white", "hamster_white", srcHamster, "eee7d8", ""),
	sourceVariant("hamster_cinnamon", "hamster", "Hamster - cinnamon", "hamster_cinnamon", srcHamster, "b56b38", ""),

	acceptedMotionVariant("ferret_sable_panda", "ferret", "Ferret - sable panda", "フェレット（セーブルパンダ）", "ferret_sable_panda", srcFerretSablePanda, srcFerretSablePandaMotion, "eee7d9", "6e5140", "Ferret", "sable panda", 2),
	acceptedMotionVariant("ferret_sable", "ferret", "Ferret - sable", "フェレット（セーブル）", "ferret_sable", srcFerretSable, srcFerretSableMotion, "8b6746", "ece0c8", "Ferret", "sable", 1),
	acceptedMotionVariant("ferret_albino", "ferret", "Ferret - albino", "フェレット（アルビノ）", "ferret_albino", srcFerretAlbino, srcFerretAlbinoMotion, "eadcc7", "c58b78", "Ferret", "albino", 2),
	shapeVariant("ferret_champagne", "ferret", "Ferret - champagne", "ferret", "c7a476", "f1e1c1"),
	acceptedMotionVariant("guinea_pig_tricolor", "guinea_pig", "Guinea pig - tricolor", "モルモット", "guinea_pig_tricolor", srcGuineaPig, srcGuineaPigMotion, "b46c32", "f0e6d2", "Guinea pig", "tricolor", 1),
	acceptedMotionVariant("guinea_pig_russian_smoke_white", "guinea_pig", "Guinea pig - Russian smoke white", "モルモット（ロシアンスモーク白）", "guinea_pig_russian_smoke_white", srcGuineaPigRussianSmoke, srcGuineaPigRussianMotion, "f1eee6", "8d8380", "Guinea pig", "Russian smoke white", 2),
	shapeVariant("guinea_pig_cream", "guinea_pig", "Guinea pig - cream", "guinea_pig", "d9b879", "f1dfb4"),
	shapeVariant("guinea_pig_black", "guinea_pig", "Guinea pig - black", "guinea_pig", "2b2926", "d8cdb9"),
	shapeVariant("hedgehog_salt_pepper", "hedgehog", "Hedgehog - salt and pepper", "hedgehog", "81786a", "e6dcc6"),
	shapeVariant("hedgehog_cinnamon", "hedgehog", "Hedgehog - cinnamon", "hedgehog", "b06e3e", "ebdcc2"),
	shapeVariant("squirrel_red", "squirrel", "Squirrel - red", "squirrel", "b45f2f", "e9c193"),
	shapeVariant("squirrel_gray", "squirrel", "Squirrel - gray", "squirrel", "85857c", "d8d4c8"),
	acceptedMotionVariant("richardsons_ground_squirrel", "ground_squirrel", "Richardson's ground squirrel", "リチャードソンジリス", "richardsons_ground_squirrel", srcGroundSquirrel, srcGroundSquirrelMotion, "a8895c", "d8c091", "Richardson's ground squirrel", "tan", 3),
	shapeVariant("fox_red", "fox", "Fox - red", "fox", "c46832", "eee2c9"),
	shapeVariant("fox_silver", "fox", "Fox - silver", "fox", "4b4d4c", "d6d3c9"),
	shapeVariant("red_panda_classic", "red_panda", "Red panda - classic", "red_panda", "b85f31", "efe1c5"),
	shapeVariant("red_panda_dark", "red_panda", "Red panda - dark", "red_panda", "6f3f2f", "e9d9bc"),
	shapeVariant("otter_brown", "otter", "Otter - brown", "otter", "6b4a31", "d4bb91"),
	acceptedMotionVariant("sugar_glider_gray", "sugar_glider", "Sugar glider", "モモンガ", "sugar_glider_gray", srcSugarGlider, srcSugarGliderMotion, "8c8d88", "f1eee4", "Sugar glider", "gray", 1),
	acceptedMotionVariant("leucistic_sugar_glider", "sugar_glider", "Leucistic sugar glider", "リューシスティックモモンガ", "leucistic_sugar_glider", srcLeucisticSugarGlider, srcLeucisticSugarGliderMotion, "f1eadf", "b69b84", "Sugar glider", "leucistic", 2),
	acceptedMotionVariant("african_dormouse", "dormouse", "African dormouse", "アフリカヤマネ", "african_dormouse", srcAfricanDormouse, srcAfricanDormouseMotion, "8d6b46", "e1c390", "African dormouse", "brown tan", 2),
	acceptedMotionVariant("american_flying_squirrel", "flying_squirrel", "American flying squirrel", "アメリカモモンガ", "american_flying_squirrel", srcAmericanFlyingSquirrel, srcAmericanFlyingSquirrelMotion, "8a7766", "d7c3a0", "American flying squirrel", "gray brown", 2),
	shapeVariant("capybara_brown", "capybara", "Capybara - brown", "capybara", "8a603d", "c79d6a"),
	shapeVariant("capybara_sand", "capybara", "Capybara - sand", "capybara", "b98a58", "e2c69c"),
	shapeVariant("tortoise_olive", "tortoise", "Tortoise - olive", "tortoise", "6f7146", "3f3d2a"),
	shapeVariant("tortoise_dark_shell", "tortoise", "Tortoise - dark shell", "tortoise", "45452f", "8b8052"),

	acceptedMotionVariant("french_bulldog_fawn", "dog", "French Bulldog - fawn", "フレンチブルドッグ", "french_bulldog_fawn", srcFrenchBulldog, srcFrenchBulldogMotion, "c49a6c", "3a3026", "French Bulldog", "fawn", 1),
	acceptedMotionVariant("toy_poodle_apricot", "dog", "Toy Poodle - apricot", "トイプードル", "toy_poodle_apricot", srcToyPoodle, srcToyPoodleMotion, "d98a3f", "5a3825", "Toy Poodle", "apricot", 1),
	acceptedMotionVariant("miniature_schnauzer_salt_pepper", "dog", "Miniature Schnauzer - salt and pepper", "ミニチュアシュナウザー", "miniature_schnauzer_salt_pepper", srcMiniatureSchnauzer, srcMiniatureSchnauzerMotion, "8e8c86", "d8d0c4", "Miniature Schnauzer", "salt and pepper", 2),
	sourceVariantMeta("labrador_yellow", "dog", "Labrador Retriever - yellow", "labrador_yellow", srcDog, "d7b46f", "", "Labrador Retriever", "yellow", 1),
	sourceVariantMeta("labrador_black", "dog", "Labrador Retriever - black", "labrador_black", srcDog, "2b2926", "", "Labrador Retriever", "black", 1),
	sourceVariantMeta("golden_retriever_golden", "dog", "Golden Retriever - golden", "golden_retriever_golden", srcDog, "c8873d", "ead9ad", "Golden Retriever", "golden", 1),
	sourceVariantMeta("german_shepherd_black_tan", "dog", "German Shepherd - black and tan", "german_shepherd_black_tan", srcDog, "3a3026", "b8874c", "German Shepherd Dog", "black and tan", 1),
	sourceVariantMeta("dachshund_red", "dog", "Dachshund - red", "dachshund_red", srcDog, "a8572b", "", "Dachshund", "red", 1),
	sourceVariantMeta("poodle_white", "dog", "Poodle - white", "poodle_white", srcDog, "eee7d8", "", "Poodle", "white", 1),
	sourceVariantMeta("beagle_tricolor", "dog", "Beagle - tricolor", "beagle_tricolor", srcDog, "7b4d2e", "f1e6d0", "Beagle", "tricolor", 1),
	sourceVariantMeta("bulldog_white_brindle", "dog", "Bulldog - white brindle", "bulldog_white_brindle", srcDog, "e7dfcf", "75513a", "Bulldog", "white brindle", 1),
	sourceVariantMeta("shiba_inu_red", "dog", "Shiba Inu - red", "shiba_inu_red", srcDog, "b7612f", "efe0c5", "Shiba Inu", "red", 2),
	sourceVariantMeta("pomeranian_orange", "dog", "Pomeranian - orange", "pomeranian_orange", srcDog, "d58233", "f2d4a2", "Pomeranian", "orange", 2),
	sourceVariantMeta("corgi_sable", "dog", "Corgi - sable", "corgi_sable", srcDog, "9b6236", "ead6b0", "Pembroke Welsh Corgi", "sable", 2),
	acceptedMotionVariant("yorkshire_terrier_longcoat", "dog", "Yorkshire Terrier - long coat", "ヨークシャーテリア", "yorkshire_terrier_longcoat", srcYorkshireTerrier, srcYorkshireTerrierMotion, "b87a3b", "242326", "Yorkshire Terrier", "long coat", 2),

	acceptedMotionVariant("maine_coon_brown_tabby", "cat", "Maine Coon - brown tabby", "メインクーン", "maine_coon_brown_tabby", srcMaineCoon, srcMaineCoonMotion, "6e5239", "a88457", "Maine Coon", "brown tabby", 1),
	acceptedMotionVariant("ragdoll_seal_bicolor", "cat", "Ragdoll - seal bicolor", "ラグドール", "ragdoll_seal_bicolor", srcRagdoll, srcRagdollMotion, "e7ddca", "594333", "Ragdoll", "seal bicolor", 1),
	sourceVariantMeta("persian_white", "cat", "Persian - white", "persian_white", srcCat, "eee8dd", "", "Persian", "white", 1),
	acceptedMotionVariant("british_shorthair_blue", "cat", "British Shorthair - blue", "ブリティッシュショートヘア", "british_shorthair_blue", srcBritishShorthair, srcBritishShorthairMotion, "737d84", "", "British Shorthair", "blue", 1),
	acceptedMotionVariant("munchkin_brown_tabby", "cat", "Munchkin - brown tabby", "マンチカン", "munchkin_brown_tabby", srcMunchkin, srcMunchkinMotion, "8b542e", "2f241c", "Munchkin", "brown tabby", 1),
	sourceVariantMeta("siamese_seal_point", "cat", "Siamese - seal point", "siamese_seal_point", srcCat, "d8c8aa", "4a3428", "Siamese", "seal point", 1),
	sourceVariantMeta("sphynx_pink", "cat", "Sphynx - pink", "sphynx_pink", srcCat, "d7a894", "", "Sphynx", "pink", 2),
	acceptedMotionVariant("scottish_fold_silver_tabby", "cat", "Scottish Fold - silver tabby", "スコティッシュフォールド", "scottish_fold_silver_tabby", srcScottishFold, srcScottishFoldMotion, "a7aaa5", "5e625f", "Scottish Fold", "silver tabby", 2),
	sourceVariantMeta("bengal_rosetted", "cat", "Bengal - rosetted", "bengal_rosetted", srcCat, "c69048", "3f2b20", "Bengal", "rosetted", 2),
	acceptedMotionVariant("domestic_shorthair_calico", "cat", "Domestic Shorthair - calico", "三毛猫", "domestic_shorthair_calico", srcDomesticShorthairCalico, srcDomesticShorthairMotion, "e8dfcd", "bd6d2e", "Domestic Shorthair", "calico", 2),
	acceptedMotionVariant("domestic_shorthair_tabby_white_stocky", "cat", "Domestic Shorthair - tabby white", "キジ白猫", "domestic_shorthair_tabby_white_stocky", srcDomesticTabbyWhite, srcDomesticTabbyWhiteMotion, "5d5145", "f1ede4", "Domestic Shorthair", "tabby and white", 2),
	sourceVariantMeta("domestic_shorthair_tuxedo", "cat", "Domestic Shorthair - tuxedo", "domestic_shorthair_tuxedo", srcCat, "242322", "eee8dd", "Domestic Shorthair", "tuxedo", 2),

	acceptedMotionVariant("holland_lop_broken_orange", "rabbit", "Holland Lop - broken orange", "ホーランドロップ", "holland_lop_broken_orange", srcHollandLop, srcHollandLopMotion, "e6d8c3", "c97833", "Holland Lop", "broken orange", 1),
	acceptedMotionVariant("netherland_dwarf_chestnut", "rabbit", "Netherland Dwarf - chestnut", "ネザーランドドワーフ", "netherland_dwarf_chestnut", srcNetherlandDwarf, srcNetherlandDwarfMotion, "8c633d", "d3b17a", "Netherland Dwarf", "chestnut", 1),
	acceptedMotionVariant("netherland_dwarf_himalayan", "rabbit", "Netherland Dwarf - Himalayan", "ネザーランドドワーフ（ヒマラヤン）", "netherland_dwarf_himalayan", srcNetherlandDwarfHimalayan, srcNetherlandDwarfHimalayanMotion, "f2e6d4", "5a4237", "Netherland Dwarf", "Himalayan", 2),
	acceptedMotionVariant("himalayan_rabbit", "rabbit", "Himalayan Rabbit", "ヒマラヤンうさぎ", "himalayan_rabbit", srcHimalayanRabbit, srcHimalayanRabbitMotion, "f3e5c8", "5b3b2d", "Himalayan rabbit", "seal point", 1),
	acceptedMotionVariant("lionhead_rabbit_brown_white", "rabbit", "Lionhead rabbit - brown white", "ライオンラビット（白茶）", "lionhead_rabbit_brown_white", srcLionheadRabbit, srcLionheadRabbitMotion, "f1ede4", "9a6845", "Lionhead rabbit", "brown white", 2),
	acceptedMotionVariant("rabbit_gray", "rabbit", "Gray rabbit", "グレーうさぎ", "rabbit_gray", srcRabbitGray, srcRabbitGrayMotion, "8f8d86", "d8d3c8", "Rabbit", "gray", 2),
	sourceVariantMeta("mini_rex_black_otter", "rabbit", "Mini Rex - black otter", "mini_rex_black_otter", srcRabbit, "2f2c29", "b89865", "Mini Rex", "black otter", 1),
	sourceVariantMeta("lionhead_tort", "rabbit", "Lionhead - tort", "lionhead_tort", srcRabbit, "a66a3a", "4f3326", "Lionhead", "tort", 2),
	sourceVariantMeta("dutch_black_white", "rabbit", "Dutch - black and white", "dutch_black_white", srcRabbit, "2c2a28", "eee7dc", "Dutch", "black and white", 2),

	acceptedMotionVariant("fancy_rat_hooded", "rat", "Fancy rat - hooded", "ファンシーラット", "fancy_rat_hooded", srcFancyRat, srcFancyRatMotion, "efe5d4", "3b332c", "Fancy rat", "hooded", 2),
	acceptedMotionVariant("fancy_rat_blue_hooded", "rat", "Fancy rat - blue hooded", "ファンシーラット（ブルーフーディッド）", "fancy_rat_blue_hooded", srcFancyRatBlueHooded, srcFancyRatBlueHoodedMotion, "d8d6cd", "6d7480", "Fancy rat", "blue hooded", 2),
	acceptedMotionVariant("fancy_rat_chocolate_self", "rat", "Fancy rat - chocolate self", "ファンシーラット（チョコレートセルフ）", "fancy_rat_chocolate_self", srcFancyRatChocolateSelf, srcFancyRatChocolateSelfMotion, "5f3e2c", "9d7350", "Fancy rat", "chocolate self", 2),
	acceptedMotionVariant("fancy_rat_cream_agouti", "rat", "Fancy rat - cream agouti", "ファンシーラット（クリームアグーチ）", "fancy_rat_cream_agouti", srcFancyRatCreamAgouti, srcFancyRatCreamAgoutiMotion, "caa06f", "ead5a7", "Fancy rat", "cream agouti", 2),
	shapeVariantMeta("fancy_mouse_white", "mouse", "Fancy mouse - white", "small_rodent", "eee7d8", "c99a8c", "Fancy mouse", "white", 2),
	shapeVariantMeta("mongolian_gerbil_agouti", "gerbil", "Mongolian gerbil - agouti", "small_rodent", "9d7448", "e5c997", "Mongolian gerbil", "agouti", 2),
	shapeVariantMeta("prairie_dog_tan", "prairie_dog", "Prairie dog - tan", "prairie_dog", "b98958", "e0c08e", "Prairie dog", "tan", 3),
	acceptedMotionVariant("chipmunk_striped", "chipmunk", "Chipmunk - striped", "シマリス", "chipmunk_striped", srcChipmunk, srcChipmunkMotion, "a06a3a", "2f241e", "Chipmunk", "striped", 3),
	acceptedMotionVariant("albino_chipmunk", "chipmunk", "Black-eyed white chipmunk", "黒目の白いホワイトシマリス", "albino_chipmunk", srcAlbinoChipmunk, srcAlbinoChipmunkMotion, "f2dfbd", "cf9a5b", "Chipmunk", "black-eyed white", 3),
	acceptedMotionVariant("true_albino_chipmunk", "chipmunk", "True albino chipmunk", "真のアルビノシマリス", "true_albino_chipmunk", srcTrueAlbinoChipmunk, srcTrueAlbinoChipmunkMotion, "f7e9d6", "d98f8a", "Chipmunk", "true albino", 3),

	shapeVariantMeta("bearded_dragon_citrus", "bearded_dragon", "Bearded dragon - citrus", "dragon", "d99b37", "f0d268", "Bearded dragon", "citrus", 2),
	sourceVariantMeta("crested_gecko_harlequin", "crested_gecko", "Crested gecko - harlequin", "crested_gecko_harlequin", srcGecko, "b77b3d", "eee0bd", "Crested gecko", "harlequin", 2),
	shapeVariantMeta("corn_snake_amelanistic", "corn_snake", "Corn snake - amelanistic", "snake", "d66b36", "f0d096", "Corn snake", "amelanistic", 2),
	shapeVariantMeta("whites_tree_frog_green", "whites_tree_frog", "White's tree frog - green", "frog", "77a95a", "d7e9b8", "White's tree frog", "green", 2),
	acceptedMotionVariant("whites_tree_frog_blue", "whites_tree_frog", "Blue White's tree frog", "水色イエアメガエル", "whites_tree_frog_blue", srcWhitesTreeFrogBlue, srcWhitesTreeFrogBlueMotion, "17b8e8", "e8dcc7", "White's tree frog", "blue", 2),
	acceptedMotionVariant("japanese_giant_salamander", "salamander", "Japanese giant salamander", "オオサンショウウオ", "japanese_giant_salamander", srcJapaneseGiantSalamander, srcJapaneseGiantSalamanderMotion, "5d4e43", "9f8f7c", "Japanese giant salamander", "mottled brown", 3),
	acceptedMotionVariant("african_fat_tailed_gecko", "gecko", "African fat-tailed gecko", "ニシアフリカトカゲモドキ", "african_fat_tailed_gecko", srcAfricanFatTailedGecko, srcAfricanFatTailedGeckoMotion, "7a5b45", "d9b07c", "African fat-tailed gecko", "brown banded", 2),

	acceptedMotionVariant("budgerigar_green_yellow", "budgerigar", "Budgerigar - green yellow", "セキセイインコ", "budgerigar_green_yellow", srcBudgerigar, srcBudgerigarMotion, "75b64c", "f0da4d", "Budgerigar", "green yellow", 1),
	acceptedMotionVariant("cockatiel_normal_gray", "cockatiel", "Cockatiel - normal gray", "オカメインコ", "cockatiel_normal_gray", srcCockatiel, srcCockatielMotion, "8d8a86", "f3d35c", "Cockatiel", "normal gray", 1),
	acceptedMotionVariant("java_sparrow_normal", "java_sparrow", "Java sparrow - normal", "文鳥", "java_sparrow_normal", srcJavaSparrow, srcJavaSparrowMotion, "8f7b83", "d44848", "Java sparrow", "normal", 1),
	acceptedMotionVariant("parrotlet_green", "parrotlet", "Parrotlet - green", "マメルリハ", "parrotlet_green", srcParrotlet, srcParrotletMotion, "80bd38", "d9e982", "Pacific parrotlet", "green", 2),
	acceptedMotionVariant("parrotlet_blue_green", "parrotlet", "Parrotlet - blue green", "マメルリハ（青緑）", "parrotlet_blue_green", srcParrotletBlueGreen, srcParrotletBlueGreenMotion, "28c2df", "136f78", "Pacific parrotlet", "blue green", 2),
	acceptedMotionVariant("lovebird_peach_faced", "lovebird", "Lovebird - peach faced", "コザクラインコ", "lovebird_peach_faced", srcLovebird, srcLovebirdMotion, "75b83e", "ef7944", "Peach-faced lovebird", "peach faced", 2),
	acceptedMotionVariant("white_wagtail", "wagtail", "White wagtail", "ハクセキレイ", "white_wagtail", srcWhiteWagtail, srcWhiteWagtailMotion, "d8d8cf", "1d1d1f", "White wagtail", "black white gray", 3),
	acceptedMotionVariant("shoebill_stork", "shoebill", "Shoebill", "ハシビロコウ", "shoebill_stork", srcShoebill, srcShoebillMotion, "7f8f91", "c8b894", "Shoebill", "gray blue", 3),
}

var runtimeVariantIDs = []string{
	"chinchilla_standard_gray",
	"chinchilla_beige",
	"chinchilla_ebony",
	"hamster_golden_syrian",
	"djungarian_hamster",
	"campbell_hamster",
	"macaroni_mouse_tan",
	"sugar_glider_gray",
	"rabbit_chestnut_agouti",
	"holland_lop_broken_orange",
	"netherland_dwarf_chestnut",
	"himalayan_rabbit",
	"gecko_gray_brown",
	"guinea_pig_tricolor",
	"fancy_rat_hooded",
	"richardsons_ground_squirrel",
	"yorkshire_terrier_longcoat",
	"chipmunk_striped",
	"true_albino_chipmunk",
	"gecko_leopard",
	"whites_tree_frog_blue",
	"budgerigar_green_yellow",
	"cockatiel_normal_gray",
	"java_sparrow_normal",
	"parrotlet_green",
	"parrotlet_blue_green",
	"lovebird_peach_faced",
	"ragdoll_seal_bicolor",
	"scottish_fold_silver_tabby",
	"french_bulldog_fawn",
	"maine_coon_brown_tabby",
	"domestic_shorthair_calico",
	"british_shorthair_blue",
	"toy_poodle_apricot",
	"munchkin_brown_tabby",
	"roborovski_hamster",
	"guinea_pig_russian_smoke_white",
	"quokka",
	"miniature_schnauzer_salt_pepper",
	"japanese_giant_salamander",
	"white_wagtail",
	"domestic_shorthair_tabby_white_stocky",
	"lionhead_rabbit_brown_white",
	"shoebill_stork",
	"leucistic_sugar_glider",
	"african_dormouse",
	"netherland_dwarf_himalayan",
	"american_flying_squirrel",
	"longhair_hamster_black_white",
	"longhair_hamster_black_white_masked",
	"djungarian_hamster_yellow",
	"djungarian_hamster_pearl_white",
	"fancy_rat_blue_hooded",
	"fancy_rat_chocolate_self",
	"fancy_rat_cream_agouti",
	"rabbit_gray",
	"african_fat_tailed_gecko",
	"ferret_sable_panda",
	"ferret_sable",
	"ferret_albino",
}

func RuntimeVariants() []Variant {
	out := make([]Variant, 0, len(runtimeVariantIDs))
	for _, id := range runtimeVariantIDs {
		if variant, ok := VariantByID(id); ok {
			out = append(out, variant)
		}
	}
	return out
}

func VariantByID(id string) (Variant, bool) {
	for _, variant := range Variants {
		if variant.ID == id {
			return variant, true
		}
	}
	return Variant{}, false
}

func deguVariant(id string, label string, spriteBase string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      "degu",
		BreedOrMorph:   "Degu",
		Color:          label,
		PopularityTier: 1,
		MotionProfile:  MotionProfileDegu,
		SourceStatus:   SourceStatusDeguMotion,
		LabelEN:        "Degu - " + label,
		LabelJA:        "Degu - " + label,
		SpriteBase:     spriteBase,
	}
}

func sourceVariant(id string, speciesID string, label string, spriteBase string, sourcePath string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      speciesID,
		BreedOrMorph:   label,
		Color:          label,
		PopularityTier: 3,
		MotionProfile:  DefaultMotionProfileForSpecies(speciesID),
		SourceStatus:   SourceStatusPrototypeOnly,
		LabelEN:        label,
		LabelJA:        label,
		SpriteBase:     spriteBase,
		SeedStage:      true,
		SourcePath:     sourcePath,
		TintHex:        tintHex,
		AccentHex:      accentHex,
	}
}

func motionSourceVariant(id string, speciesID string, label string, spriteBase string, sourcePath string, motionSourcePath string, tintHex string, accentHex string) Variant {
	variant := sourceVariant(id, speciesID, label, spriteBase, sourcePath, tintHex, accentHex)
	variant.SourceStatus = SourceStatusMotionDraft
	variant.MotionSourcePath = motionSourcePath
	return variant
}

func acceptedMotionVariant(id string, speciesID string, labelEN string, labelJA string, spriteBase string, sourcePath string, motionSourcePath string, tintHex string, accentHex string, breedOrMorph string, color string, popularityTier int) Variant {
	variant := sourceVariantMeta(id, speciesID, labelEN, spriteBase, sourcePath, tintHex, accentHex, breedOrMorph, color, popularityTier)
	variant.SourceStatus = SourceStatusMotionAccepted
	variant.MotionSourcePath = motionSourcePath
	variant.LabelJA = labelJA
	return variant
}

func sourceVariantMeta(id string, speciesID string, label string, spriteBase string, sourcePath string, tintHex string, accentHex string, breedOrMorph string, color string, popularityTier int) Variant {
	variant := sourceVariant(id, speciesID, label, spriteBase, sourcePath, tintHex, accentHex)
	variant.BreedOrMorph = breedOrMorph
	variant.Color = color
	variant.PopularityTier = popularityTier
	return variant
}

func shapeVariant(id string, speciesID string, label string, shape string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      speciesID,
		BreedOrMorph:   label,
		Color:          label,
		PopularityTier: 3,
		MotionProfile:  DefaultMotionProfileForSpecies(speciesID),
		SourceStatus:   SourceStatusPrototypeOnly,
		LabelEN:        label,
		LabelJA:        label,
		SpriteBase:     id,
		SeedStage:      true,
		Shape:          shape,
		TintHex:        tintHex,
		AccentHex:      accentHex,
	}
}

func shapeVariantMeta(id string, speciesID string, label string, shape string, tintHex string, accentHex string, breedOrMorph string, color string, popularityTier int) Variant {
	variant := shapeVariant(id, speciesID, label, shape, tintHex, accentHex)
	variant.BreedOrMorph = breedOrMorph
	variant.Color = color
	variant.PopularityTier = popularityTier
	return variant
}

func DeguVariantIDs() []string {
	ids := make([]string, 0, len(Variants))
	for _, variant := range Variants {
		if variant.SpeciesID == "degu" {
			ids = append(ids, variant.ID)
		}
	}
	return ids
}

func SeedVariants() []Variant {
	out := make([]Variant, 0)
	for _, variant := range Variants {
		if variant.SeedStage {
			out = append(out, variant)
		}
	}
	return out
}

func SpeciesByID(id string) (Species, bool) {
	for _, species := range SpeciesList {
		if species.ID == id {
			return species, true
		}
	}
	return Species{}, false
}

func ProfileForSpecies(id string) string {
	if species, ok := SpeciesByID(id); ok {
		return species.Profile
	}
	return ""
}

func DefaultMotionProfileForSpecies(speciesID string) string {
	switch speciesID {
	case "degu":
		return MotionProfileDegu
	case "rabbit", "quokka":
		return MotionProfileRabbitHop
	case "dog":
		return MotionProfileDogTrot
	case "cat":
		return MotionProfileCatStalk
	case "gecko", "crested_gecko", "salamander":
		return MotionProfileGeckoCrawl
	case "tortoise":
		return MotionProfileTortoisePlod
	case "ferret":
		return MotionProfileFerretSlink
	case "guinea_pig":
		return MotionProfileGuineaPigWaddle
	case "hedgehog":
		return MotionProfileHedgehogShuffle
	case "squirrel", "flying_squirrel", "ground_squirrel", "chipmunk":
		return MotionProfileSquirrelBound
	case "fox":
		return MotionProfileFoxTrot
	case "red_panda":
		return MotionProfileRedPandaAmble
	case "otter":
		return MotionProfileOtterSlide
	case "sugar_glider":
		return MotionProfileSugarGliderSkitter
	case "capybara":
		return MotionProfileCapybaraLumber
	case "corn_snake":
		return MotionProfileSnakeSlither
	case "bearded_dragon":
		return MotionProfileDragonPlod
	case "whites_tree_frog":
		return MotionProfileFrogHop
	case "budgerigar", "cockatiel", "java_sparrow", "parrotlet", "lovebird", "wagtail", "shoebill":
		return MotionProfileBirdHop
	default:
		return MotionProfileSmallRodentScurry
	}
}

func MotionProfileForVariant(variant Variant) string {
	if variant.MotionProfile != "" {
		return variant.MotionProfile
	}
	return DefaultMotionProfileForSpecies(variant.SpeciesID)
}

func WheelCapableSpecies(speciesID string) bool {
	switch speciesID {
	case "chinchilla", "hamster":
		return true
	default:
		return false
	}
}

func WheelCapableMotionProfile(profile string) bool {
	return false
}

func WheelCapableVariant(variant Variant) bool {
	switch variant.ID {
	case "chinchilla_standard_gray", "chinchilla_beige", "chinchilla_ebony", "hamster_golden_syrian":
		return true
	default:
		return false
	}
}
