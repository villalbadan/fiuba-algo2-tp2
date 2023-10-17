package gram

import (
	abb "algogram/abb"
	"strings"
)

type publicacionImp struct {
	contenido  string
	id         int
	publicador Usuario
	likes      abb.DiccionarioOrdenado[string, int]
}

func mayorEntreStrings(clave1, clave2 string) int {
	return strings.Compare(clave1, clave2)
}

func CrearPublicacion(mensaje string, emisor Usuario, idActual int) Publicacion {

	nuevoPost := new(publicacionImp)
	nuevoPost.contenido = mensaje
	nuevoPost.id = idActual
	nuevoPost.publicador = emisor
	nuevoPost.likes = abb.CrearABB[string, int](mayorEntreStrings)

	return nuevoPost
}

func (publicacion publicacionImp) TextoPublicacion() string {
	return publicacion.contenido
}
func (publicacion publicacionImp) Id() int {
	return publicacion.id
}

func (publicacion publicacionImp) VerPublicador() Usuario {
	return publicacion.publicador
}

func (publicacion *publicacionImp) Likear(usuario Usuario) {
	publicacion.likes.Guardar(usuario.Nombre(), 0)
}

func (publicacion *publicacionImp) MostrarLikes() []string {
	listaUsuarios := make([]string, 0, 1)

	publicacion.likes.Iterar(func(clave string, dato int) bool {
		listaUsuarios = append(listaUsuarios, clave)
		return true
	})

	return listaUsuarios
}

func (publicacion publicacionImp) CantLikes() int {
	return publicacion.likes.Cantidad()
}
