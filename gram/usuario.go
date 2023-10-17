package gram

type Usuario interface {

	// Id devuelve un int como id del usuario
	Id() int

	// Nombre devuelve el dato asociado al nombre del usuario
	Nombre() string

	// VerAfinidad devuelve la afinidad entre el usuario actualmente logueado y el usuario ingresado
	VerAfinidad(usuarioAComparar Usuario) int

	// ActualizarFeed agrega un puntero a una Publicacion en el feed del usuario actual
	ActualizarFeed(post *Publicacion)

	// ProximoPost devuelve el próximo post en prioridad. Si el feed esta vacio, devuelve nil
	// y el error con el mensaje "No hay posts que mostrar"
	ProximoPost() (*Publicacion, error)
}
