package gram

import (
	heap "algogram/heap"
	"errors"
)

type usuarioImp struct {
	nombre        string
	id            int
	prioridadFeed heap.ColaPrioridad[*Publicacion]
}

func CrearUsuario(nombre string, id int) Usuario {
	nuevoUsuario := new(usuarioImp)
	nuevoUsuario.nombre = nombre
	nuevoUsuario.id = id
	nuevoUsuario.prioridadFeed = heap.CrearHeap[*Publicacion](nuevoUsuario.cmpAfinidad)
	return nuevoUsuario
}

// Auxiliares calculo de afinidad ---------------------------------------------------------------------

// abs devuelve el valor absoluto de un entero
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func mayorEntreInts(clave1, clave2 int) int {
	return clave1 - clave2
}

// cmpAfinidad Devuelve un int con la diferencia de afinidad entre los dos posteos pasados como parametro
func (usuario *usuarioImp) cmpAfinidad(post1, post2 *Publicacion) int {

	// afinidad entre publicadores
	mayorAfinidad := mayorEntreInts(usuario.VerAfinidad((*post2).VerPublicador()), usuario.VerAfinidad((*post1).VerPublicador()))
	if mayorAfinidad != 0 {
		return mayorAfinidad
	}

	// post mas reciente
	return mayorEntreInts((*post2).Id(), (*post1).Id())
}

// PRIMITIVAS USUARIO ----------------------------------------------------------------------------------------

func (usuario usuarioImp) Nombre() string {
	return usuario.nombre
}

func (usuario usuarioImp) Id() int {
	return usuario.id
}

func (usuario usuarioImp) VerAfinidad(usuario2 Usuario) int {
	return abs(usuario.id - usuario2.Id())
}

func (usuario *usuarioImp) ActualizarFeed(post *Publicacion) {
	usuario.prioridadFeed.Encolar(post)
}

func (usuario *usuarioImp) ProximoPost() (*Publicacion, error) {
	if usuario.prioridadFeed.EstaVacia() {
		return nil, errors.New("No hay posts que mostrar")
	}

	return usuario.prioridadFeed.Desencolar(), nil
}
