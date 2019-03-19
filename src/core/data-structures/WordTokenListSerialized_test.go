package data_structures

/*
func TestWordTokenListSerialized_Deserialize(t *testing.T) {

	var wordTokenListBeforeSerialization *WordTokenList
	var wordTokenListSerialized WordTokenListSerialized
	var digestBeforeSerialization string
	var digestAfterSerialization string
	var outputFromRead interface{}


	wordTokenListBeforeSerialization = BuildWordTokenList()

	(*wordTokenListBeforeSerialization).InsertWordToken(BuildWordToken("Andrea", 5))
	(*wordTokenListBeforeSerialization).InsertWordToken(BuildWordToken("Graziani", 5))
	(*wordTokenListBeforeSerialization).InsertWordToken(BuildWordToken("Diana", 5))

	wordTokenListSerialized = (*wordTokenListBeforeSerialization).Serialize()
	digestBeforeSerialization, _ = utility.SHA512(wordTokenListSerialized)

	_ = utility.WriteToLocalDisk(wordTokenListSerialized)

	dd := WordTokenListSerialized{}

	outputFromRead, _ = utility.ReadFromLocalDisk(digestBeforeSerialization, &dd)
	digestAfterSerialization, _ = utility.SHA512(outputFromRead)



	(*(outputFromRead).(WordTokenListSerialized).Deserialize()).Print()

	if strings.Compare(digestBeforeSerialization, digestAfterSerialization) != 0 {
		log.Fatal("Digest NOT correct!")
	}

}
*/
