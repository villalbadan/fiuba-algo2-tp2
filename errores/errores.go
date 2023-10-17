package errores

type ErrorLeerArchivo struct{}

func (e ErrorLeerArchivo) Error() string {
	return "ERROR: Lectura de archivos"
}

type ErrorParametros struct{}

func (e ErrorParametros) Error() string {
	return "ERROR: Faltan parámetros"
}

type UsuarioInexistenteError struct{}

func (e UsuarioInexistenteError) Error() string {
	return "Error: usuario no existente"
}

type LoginYaRealizadoError struct{}

func (e LoginYaRealizadoError) Error() string {
	return "Error: Ya habia un usuario loggeado"
}

type SinLoggearError struct{}

func (e SinLoggearError) Error() string {
	return "Error: no habia usuario loggeado"
}

type VerSiguienteError struct {
}

func (e VerSiguienteError) Error() string {
	return "Usuario no loggeado o no hay mas posts para ver"
}

type LikearError struct{}

func (e LikearError) Error() string {
	return "Error: Usuario no loggeado o Post inexistente"
}

type MostrarLikesError struct{}

func (e MostrarLikesError) Error() string {
	return "Error: Post inexistente o sin likes"
}

type ErrorComando struct{}

func (e ErrorComando) Error() string {
	return "Error: comando inexistente"
}
