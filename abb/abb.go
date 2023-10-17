package diccionario

import (
	TDAPila "algogram/abb/pila"
)

const (
	VALOR_CMP = 0
)

// ################################### ESTRUCTURAS ###############################################################

type funcCmp[K comparable] func(K, K) int

type ab[K comparable, V any] struct {
	raiz     *nodoAb[K, V]
	cantidad int
	cmp      funcCmp[K]
}

type nodoAb[K comparable, V any] struct {
	izq   *nodoAb[K, V]
	der   *nodoAb[K, V]
	clave K
	dato  V
}

type iteradorDict[K comparable, V any] struct {
	diccionario   *ab[K, V]
	actual        *nodoAb[K, V]
	rangoMin      *K
	rangoMax      *K
	pilaElementos TDAPila.Pila[*nodoAb[K, V]]
}

// ##############################################################################################################

func CrearABB[K comparable, V any](funcion_cmp funcCmp[K]) DiccionarioOrdenado[K, V] {
	dict := new(ab[K, V])
	dict.cmp = funcion_cmp
	return dict
}

func (dict *ab[K, V]) buscar(clave K, nodoActual **nodoAb[K, V]) **nodoAb[K, V] {
	if (*nodoActual) == nil {
		return nodoActual
	}

	comparacion := dict.cmp(clave, (*nodoActual).clave)
	// clave a evaluar es menor a la clave actual
	if comparacion < VALOR_CMP {
		return dict.buscar(clave, &(*nodoActual).izq)
	}

	// clave a evaluar es mayor a la clave actual
	if comparacion > VALOR_CMP {
		return dict.buscar(clave, &(*nodoActual).der)
	}

	// clave a evaluar es la clave del nodo
	return nodoActual

}

// ################################### Aux. Borrar #########################################################

// noTieneHijos recibe un nodo, y devuelve True si su hijo izquierdo y su hijo derecho son ambos nil
func noTieneHijos[K comparable, V any](nodo **nodoAb[K, V]) bool {
	return (*nodo).izq == nil && (*nodo).der == nil
}

// reemplazante recibe un nodo de un abb y devuelve el nodo menor hacia la izquierda
// (== con un solo hijo a la derecha o sin hijos) de toda la rama, incluyendose a si mismo
func (dict *ab[K, V]) reemplazante(nodoActual **nodoAb[K, V]) *nodoAb[K, V] {

	if noTieneHijos(nodoActual) || (*nodoActual).izq == nil {
		return *nodoActual
	}
	return dict.reemplazante(&(*nodoActual).izq)
}

// transplantar recibe un nodo a borrar (con hijos). Si tiene un solo hijo, lo reemplaza por su hijo,
// si tiene dos hijos, busca el reemplazante en la rama de su hijo derecho.
func (dict *ab[K, V]) transplantar(nodo **nodoAb[K, V]) {

	// nodo con un solo hijo
	if (*nodo).izq == nil && (*nodo).der != nil {
		*nodo = (*nodo).der
		return
	}

	if (*nodo).izq != nil && (*nodo).der == nil {
		*nodo = (*nodo).izq
		return
	}

	// nodo con dos hijos
	// Busco reemplazante menor de la derecha
	nuevoNodo := dict.reemplazante(&(*nodo).der)
	nuevaClave := nuevoNodo.clave
	nuevoDato := dict.Borrar(nuevaClave)
	dict.cantidad++ // Para contrarrestar el Borrar de la linea de arriba

	// piso datos
	(*nodo).clave = nuevaClave
	(*nodo).dato = nuevoDato
}

// ################################### PRIMITIVAS DICCIONARIO #################################################

func (dict *ab[K, V]) Guardar(clave K, dato V) {
	nodo := dict.buscar(clave, &dict.raiz)
	if *nodo != nil {
		(*nodo).dato = dato
		return
	}

	*nodo = &nodoAb[K, V]{clave: clave, dato: dato}
	dict.cantidad++
}

func (dict *ab[K, V]) Pertenece(clave K) bool {
	return *(dict.buscar(clave, &dict.raiz)) != nil
}

func (dict *ab[K, V]) Obtener(clave K) V {
	nodo := dict.buscar(clave, &dict.raiz)
	if *nodo == nil {
		panic("La clave no pertenece al diccionario")
	}
	return (*nodo).dato
}

func (dict *ab[K, V]) Cantidad() int {
	return dict.cantidad
}

func (dict *ab[K, V]) Borrar(clave K) V {
	nodo := dict.buscar(clave, &dict.raiz)
	if *nodo == nil {
		panic("La clave no pertenece al diccionario")
	}

	// dato del nodo a borrar
	borrado := (*nodo).dato
	dict.cantidad--

	if noTieneHijos(nodo) {
		// nodo sin hijos
		*nodo = nil
	} else {
		// nodo con hijos
		dict.transplantar(nodo)
	}

	return borrado

}

// Iterador interno ------------------------------------------------->>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func (dict ab[K, V]) Iterar(visitar func(K, V) bool) {
	continuo := true
	continuoPtr := &continuo
	dict.raiz.iterar(visitar, continuoPtr)
}

func (nodo *nodoAb[K, V]) iterar(visitar func(K, V) bool, continuoPtr *bool) {

	if nodo == nil {
		return
	}

	if nodo.izq != nil {
		nodo.izq.iterar(visitar, continuoPtr)
	}
	if *continuoPtr && !visitar(nodo.clave, nodo.dato) {
		*continuoPtr = false
		return
	}
	if nodo.der != nil {
		nodo.der.iterar(visitar, continuoPtr)
	}
	return

}

// ################################### PRIMITIVAS ITERADOR EXTERNO ################################################
func (dict *ab[K, V]) crearIter(desde *K, hasta *K) IterDiccionario[K, V] {
	iter := iteradorDict[K, V]{diccionario: dict, rangoMin: desde, rangoMax: hasta}
	iter.pilaElementos = TDAPila.CrearPilaDinamica[*nodoAb[K, V]]()

	if hasta != nil && desde != nil && dict.cmp(*desde, *hasta) > VALOR_CMP {
		return &iter
	}

	if desde == nil {
		iter.actual = dict.raiz.buscarHijosIzquierdayApilar(iter.pilaElementos)
	} else {
		iter.actual = dict.raiz.buscarMinimo(iter.pilaElementos, dict.cmp, iter.rangoMin)
	}
	return &iter
}

// buscarHijosIzquierdayApilar dado un nodo, busca todos los hijos a la izquierda y los apila
// se detiene cuando no hay más hijos a la izquierda que apilar
func (nodo *nodoAb[K, V]) buscarHijosIzquierdayApilar(pila TDAPila.Pila[*nodoAb[K, V]]) *nodoAb[K, V] {
	if nodo == nil {
		return nil
	}
	pila.Apilar(nodo)
	if nodo.izq == nil {
		return nodo
	}
	return nodo.izq.buscarHijosIzquierdayApilar(pila)
}

func (dict *ab[K, V]) Iterador() IterDiccionario[K, V] {
	return dict.crearIter(nil, nil)
}

func (iter *iteradorDict[K, V]) HaySiguiente() bool {
	if iter.rangoMax != nil && iter.actual != nil {
		resultadoCmp := iter.diccionario.cmp(iter.actual.clave, *iter.rangoMax)
		return !iter.pilaElementos.EstaVacia() && resultadoCmp <= VALOR_CMP
	}
	return !iter.pilaElementos.EstaVacia()
}

func (iter *iteradorDict[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	clave, dato := iter.pilaElementos.VerTope().clave, iter.pilaElementos.VerTope().dato
	return clave, dato
}

func (iter *iteradorDict[K, V]) Siguiente() K {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	nodoActual := iter.pilaElementos.Desapilar()
	if nodoActual.der != nil {
		nodoActual.der.buscarHijosIzquierdayApilar(iter.pilaElementos)
	}
	if !iter.pilaElementos.EstaVacia() {
		iter.actual = iter.pilaElementos.VerTope()
	}
	return nodoActual.clave
}

// ################################### PRIMITIVAS DICCIONARIO ORDENADO #############################################

func (dict ab[K, V]) IteradorRango(desde *K, hasta *K) IterDiccionario[K, V] {
	return dict.crearIter(desde, hasta)

}

func (dict ab[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	continuo := true
	continuoPtr := &continuo

	if desde == nil && hasta == nil {
		dict.raiz.iterar(visitar, continuoPtr)
		return
	}

	if desde == nil || hasta == nil || dict.cmp(*hasta, *desde) > VALOR_CMP {
		dict.raiz.iterarRango(desde, hasta, visitar, dict.cmp, continuoPtr)
	}
}

func (nodo *nodoAb[K, V]) iterarRango(desde *K, hasta *K, visitar func(K, V) bool, cmp funcCmp[K], continuoPtr *bool) {

	if nodo == nil {
		return
	}

	// CONDICIONES de rango
	var (
		mayorADesde       = desde != nil && cmp(nodo.clave, *desde) > VALOR_CMP
		menorAHasta       = hasta != nil && cmp(nodo.clave, *hasta) < VALOR_CMP
		mayorOIgualADesde = desde != nil && cmp(nodo.clave, *desde) >= VALOR_CMP
		menorOIgualAHasta = hasta != nil && cmp(nodo.clave, *hasta) <= VALOR_CMP
		soloHasta         = desde == nil && menorOIgualAHasta
		soloDesde         = mayorOIgualADesde && hasta == nil
	)

	if (desde == nil) || mayorADesde {
		nodo.izq.iterarRango(desde, hasta, visitar, cmp, continuoPtr)
	}

	if *continuoPtr &&
		(soloHasta || soloDesde || (mayorOIgualADesde && menorOIgualAHasta)) && !visitar(nodo.clave, nodo.dato) {
		*continuoPtr = false
		return
	}

	if (hasta == nil) || menorAHasta {
		nodo.der.iterarRango(desde, hasta, visitar, cmp, continuoPtr)
	}

}

func (nodo *nodoAb[K, V]) buscarMinimo(pila TDAPila.Pila[*nodoAb[K, V]], cmp funcCmp[K], desde *K) *nodoAb[K, V] {
	if nodo == nil {
		return nil
	}

	comparacion := cmp(*desde, nodo.clave)
	// la clave inicial es menor a la clave del nodo en el que estamos parados
	// la siguiente clave a evaluar es mayor o igual
	if comparacion < VALOR_CMP && nodo.izq != nil {
		pila.Apilar(nodo)
		return nodo.izq.buscarMinimo(pila, cmp, desde)
	}

	// la clave inicial es mayor a la clave del nodo en el que estamos parados
	if comparacion > VALOR_CMP && nodo.der != nil {
		return nodo.der.buscarMinimo(pila, cmp, desde)
	}

	// igual o más proximo a la clave actual dentro del rango
	if comparacion <= VALOR_CMP {
		pila.Apilar(nodo)
	}
	return nodo

}
