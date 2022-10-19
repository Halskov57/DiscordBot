package emotes

type Emote struct {
	Emote string `bson:"emote"`
	Count int    `bson:"count"`
	User  string `bson:"user"`
	Guild string `bson:"guild"`
}

// I need to add 4 function wrappers here that extracts the 4 fields of the emote

