package gram

type Publicacion interface {

	// TextoPublicacion devuelve el texto de la publicación
	TextoPublicacion() string

	// Id devuelve un int como id de la publicación, el cual es otorgado en orden de creación
	Id() int

	// VerPublicador devuelve el nombre de usuario del publicador del posteo
	VerPublicador() Usuario

	// Likear agrega el usuario indicado a los usuarios que likearon la publicacion
	Likear(usuario Usuario)

	// MostrarLikes devuelve una lista con los nombres de usuario que dieron like a la publicacion
	MostrarLikes() []string

	// CantLikes devuelve la cantidad de likes que tiene la publicación
	CantLikes() int
}
