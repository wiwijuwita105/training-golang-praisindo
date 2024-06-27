package config

type ProfileRiskCategory string

const (
	ProfileRiskCategoryConservative ProfileRiskCategory = "Conservative"
	ProfileRiskCategoryModerate     ProfileRiskCategory = "Moderate"
	ProfileRiskCategoryBalanced     ProfileRiskCategory = "Balanced"
	ProfileRiskCategoryGrowth       ProfileRiskCategory = "Growth"
	ProfileRiskCategoryAggresive    ProfileRiskCategory = "Aggresive"
)

type ProfileRisk struct {
	MinScore   int
	MaxScore   int
	Category   ProfileRiskCategory
	Definition string
}

var RiskMapping = []ProfileRisk{
	{
		MinScore: 0,
		MaxScore: 11,
		Category: ProfileRiskCategoryConservative,
		Definition: "Tujuan utama Anda adalah untuk melindungi modal/dana yang ditempatkan dan Anda tidak memiliki toleransi " +
			"sama sekali terhadap perubahan harga/nilai dari dana investasinya tersebut. " +
			"Anda memiliki pengalaman yang sangat terbatas atau tidak memiliki pengalaman sama sekali mengenai produk investasi.",
	},
	{
		MinScore:   12,
		MaxScore:   19,
		Category:   ProfileRiskCategoryModerate,
		Definition: "Anda memiliki toleransi yang rendah dengan perubahan harga/nilai dari dana investasi dan risiko investasi.",
	},
	{
		MinScore: 20,
		MaxScore: 28,
		Category: ProfileRiskCategoryBalanced,
		Definition: "Anda memiliki toleransi yang cukup terhadap produk investasi dan dapat menerima perubahan yang besar dari " +
			"harga/nilai dari harga yang diinvestasikan.",
	},
	{
		MinScore: 29,
		MaxScore: 35,
		Category: ProfileRiskCategoryGrowth,
		Definition: "Anda memiliki toleransi yang cukup tinggi dan dapat menerima perubahan yang besar dari harga/nilai portfolio" +
			"pada produk investasi yang diinvestasikan." +
			"Pada umumnya Anda sudah pernah atau berpengalaman dalam berinvestasi di produk investasi.",
	},
	{
		MinScore: 36,
		MaxScore: 40,
		Category: ProfileRiskCategoryAggresive,
		Definition: "Anda sangat berpengalaman terhadap produk investasi dan memiliki toleransi yang sangat tinggi atas" +
			"produk-produk investasi. Anda bahkan dapat menerima perubahan signifikan pada modal/nilai investasi." +
			"Pada umumnya portfolio Anda sebagian besar dialokasikan pada produk investasi.",
	},
}
