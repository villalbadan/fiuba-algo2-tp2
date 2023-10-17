package main

import (
	"algogram/errores"
	gram "algogram/gram"
	hash "algogram/hash"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ARCHIVOS_INICIO    = 1
	COMANDO            = 0
	INIT_ARRAYS        = 50
	FACTOR_REDIMENSION = 2
)

// ############### ---------------------------------------------------------------------------------------------------

func usuarioLoggeado(actual gram.Usuario) bool {
	return actual != nil
}

// MOSTRAR LIKES ----------------------------------------------------------------------------------------------------

// mostrarLikes imprime por pantalla la cantidad de likes que tiene el post, asi como los nombres de los usuarios
// que le dieron like. Si el post no tiene likes o el id es incorrecto (en formato o post no existe), devuelve
// "Error: Post inexistente o sin likes"
func mostrarLikes(id string, feed []*gram.Publicacion, cantPosts int) {
	idInt, err := strconv.Atoi(id)
	if idInt >= cantPosts || err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.MostrarLikesError{})
		return
	}

	// repeticion incomoda de codigo, pero Go no me desreferencia el puntero dentro del if
	// y en el caso de haber algun error en el id, el programa se romperia si no se chequea antes de llamar a CantLikes
	publi := *(feed[idInt])
	if publi.CantLikes() == 0 {
		fmt.Fprintf(os.Stdout, "%s\n", errores.MostrarLikesError{})
		return
	}

	likes := publi.MostrarLikes()
	fmt.Fprintf(os.Stdout, "El post tiene %d likes:\n", publi.CantLikes())
	for i := range likes {
		fmt.Fprintf(os.Stdout, "\t%s\n", likes[i])
	}
}

// LIKEAR -----------------------------------------------------------------------------------------------------------

// likear agrega el usuario actual a la lista de usuarios que likearon el posteo con el id indicado.
// Si no hay un usuario loggeado o no existe un post con ese id, devuelve "Error: Usuario no loggeado o Post inexistente"
func likear(id string, actual gram.Usuario, feed []*gram.Publicacion, cantPosts int) {
	idInt, err := strconv.Atoi(id)
	if !usuarioLoggeado(actual) || idInt > cantPosts-1 || err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.LikearError{})
		return
	}

	publi := *(feed[idInt])
	publi.Likear(actual)
	fmt.Fprintf(os.Stdout, "Post likeado\n")
}

// VER SIGUIENTE ----------------------------------------------------------------------------------------------------

// verSiguiente muestra el siguiente post en prioridad segun la finidad entre usuarios.
// Si no hay un uusario loggeado o no hay más posts que ver, devuelve "Usuario no loggeado o no hay mas posts para ver".
func verSiguiente(actual gram.Usuario) {

	if !usuarioLoggeado(actual) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.VerSiguienteError{})
		return
	}

	sigPost, err := actual.ProximoPost()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.VerSiguienteError{})
		return
	}

	fmt.Fprintf(os.Stdout, "Post ID %d\n", (*sigPost).Id())
	fmt.Fprintf(os.Stdout, "%s dijo: %s\n", (*sigPost).VerPublicador().Nombre(), (*sigPost).TextoPublicacion())
	fmt.Fprintf(os.Stdout, "Likes: %d\n", (*sigPost).CantLikes())
}

// PUBLICAR ---------------------------------------------------------------------------------------------------------

// actualizarFeeds agrega la publicacion a cada usario existente en el hash de usuarios
func actualizarFeeds(usuariosHash hash.Diccionario[string, gram.Usuario], nuevaPubli *gram.Publicacion) {
	usuariosHash.Iterar(func(clave string, dato gram.Usuario) bool {
		if clave != (*nuevaPubli).VerPublicador().Nombre() {
			dato.ActualizarFeed(nuevaPubli)
		}
		return true
	})
}

func redimensionar(nuevoTamanio int, feed []*gram.Publicacion) []*gram.Publicacion {
	temp := feed
	feed = make([]*gram.Publicacion, nuevoTamanio)
	copy(feed, temp)

	return feed
}

func guardarEnFeed(nuevo *gram.Publicacion, feed []*gram.Publicacion, cant *int) []*gram.Publicacion {
	// controlar el estado del feed, si se necesita redimensionar
	// no se contempla redimensionar para abajo ya que los posteos no se pueden borrar
	if len(feed) <= *cant {
		feed = redimensionar(*cant*FACTOR_REDIMENSION, feed)
	}

	// guardar en array de feeds
	feed[*cant] = nuevo
	*cant += 1

	return feed

}

func publicar(texto []string, actual gram.Usuario, feed []*gram.Publicacion, usuariosHash hash.Diccionario[string, gram.Usuario], cantPosts *int) []*gram.Publicacion {
	if !usuarioLoggeado(actual) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.SinLoggearError{})
		return feed
	}

	// convertimos de slice a string
	textoCadena := strings.Join(texto, " ")

	// creamos la publicacion y actualizamos los feeds
	nuevaPubli := gram.CrearPublicacion(textoCadena, actual, *cantPosts)
	feed = guardarEnFeed(&nuevaPubli, feed, cantPosts)
	actualizarFeeds(usuariosHash, &nuevaPubli)
	fmt.Fprintf(os.Stdout, "Post publicado\n")
	return feed
}

// DESLOGGEAR -------------------------------------------------------------------------------------------------------

func desloggear(actual gram.Usuario, usuarios hash.Diccionario[string, gram.Usuario]) {
	if !usuarioLoggeado(actual) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.SinLoggearError{})
		return
	}

	usuarios.Guardar(actual.Nombre(), actual)
	fmt.Fprintf(os.Stdout, "Adios\n")
}

// LOGGEAR ----------------------------------------------------------------------------------------------------------

// existe devuelve True si existe el usuario pasado por parametro en el hash de usuarios
func existe(nombre string, usuariosHash hash.Diccionario[string, gram.Usuario]) bool {
	return usuariosHash.Pertenece(nombre)
}

func loggear(ingresado []string, usuariosHash hash.Diccionario[string, gram.Usuario], actual gram.Usuario) gram.Usuario {
	nombre := strings.Join(ingresado, " ")

	if usuarioLoggeado(actual) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.LoginYaRealizadoError{})
		return actual
	}

	if !existe(nombre, usuariosHash) {
		fmt.Fprintf(os.Stdout, "%s\n", errores.UsuarioInexistenteError{})
		return actual
	}

	fmt.Fprintf(os.Stdout, "Hola %s\n", nombre)
	return usuariosHash.Obtener(nombre)

}

// INICIO DEL PROGRAMA -----------------------------------------------------------------------------------------------

// leerUsuarios recibe la ruta del archivo con la lista de usuarios y lo devuelve como una lista de strings.
// Si hay algún problema al leer el archivo, imprime "ERROR: Lectura de archivos"
func leerUsuarios(rutaArchivo string) []string {
	temp := make([]string, 0, INIT_ARRAYS)
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorLeerArchivo{})
	}
	defer archivo.Close()

	s := bufio.NewScanner(archivo)
	for s.Scan() {
		temp = append(temp, s.Text())
	}
	err = s.Err()
	if err != nil {
		fmt.Println(err)
	}

	return temp

}

// prepararUsuarios recibe la ruta del archivo con nombres de usuario y devuelve un hash de usuarios
func prepararUsuarios(archivoUsuarios string) hash.Diccionario[string, gram.Usuario] {

	// leer archivo
	nombres := leerUsuarios(archivoUsuarios)
	usuariosHash := hash.CrearHash[string, gram.Usuario]()
	for i := range nombres {
		nuevoUsuario := gram.CrearUsuario(nombres[i], i)
		usuariosHash.Guardar(nombres[i], nuevoUsuario)
	}

	return usuariosHash

}

func inicializar(args []string) bool {

	// nro de archivos correcto
	if len(args) != ARCHIVOS_INICIO {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorParametros{})
		return false
	}

	// archivo existe
	_, err := os.Stat(args[0])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorLeerArchivo{})
		return false
	}
	return true
}

// ############### ---------------------------------------------------------------------------------------------------

func main() {

	argumentos := os.Args
	// argumentos := []string{"algogram.exe", "/home/dani/GolandProjects/algo2/Tp2/03_usuarios"}
	if inicializar(argumentos[1:]) {
		// argumentos estan correctos asi que preparamos las variables a usar
		var (
			usuarios  = prepararUsuarios(argumentos[1])
			actual    gram.Usuario
			feed      = make([]*gram.Publicacion, INIT_ARRAYS)
			cantPosts int
		)

		// lectura stdin
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			args := strings.Split(s.Text(), " ")
			switch args[COMANDO] {

			case "login":
				actual = loggear(args[COMANDO+1:], usuarios, actual)
			case "logout":
				desloggear(actual, usuarios)
				actual = nil
			case "publicar":
				feed = publicar(args[COMANDO+1:], actual, feed, usuarios, &cantPosts)
			case "ver_siguiente_feed":
				verSiguiente(actual)
			case "likear_post":
				likear(args[COMANDO+1], actual, feed, cantPosts)
			case "mostrar_likes":
				mostrarLikes(args[COMANDO+1], feed, cantPosts)
			default:
				fmt.Fprintf(os.Stdout, "%s\n", errores.ErrorComando{})

			}
		}
	}
}
