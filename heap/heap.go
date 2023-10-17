package cola_prioridad

const (
	_TAMANIO_INICIAL    = 10
	_FACTOR_REDIMENSION = 2
	_PROPORCION_MIN     = 4
)

type funcCmp[T comparable] func(T, T) int

type colaPrioridad[T comparable] struct {
	datos    []T
	cantidad int
	cmp      funcCmp[T]
}

func CrearHeap[T comparable](funcion_cmp func(T, T) int) ColaPrioridad[T] {
	heap := new(colaPrioridad[T])
	heap.datos = make([]T, _TAMANIO_INICIAL)
	heap.cmp = funcion_cmp
	return heap
}

// HEAPIFY -----------------------------------------------------------------------------------------------------------

func (heap *colaPrioridad[T]) heapify(arreglo []T, funcion_cmp func(T, T) int) {
	heap.cmp = funcion_cmp
	for i := heap.cantidad - 1; i >= 0; i-- {
		heap.downheap(heap.datos[i], i)
	}
}

func CrearHeapArr[T comparable](arreglo []T, funcion_cmp func(T, T) int) ColaPrioridad[T] {
	heapArr := new(colaPrioridad[T])

	if len(arreglo) == 0 {
		heapArr.datos = make([]T, _TAMANIO_INICIAL)
	} else {
		heapArr.datos = make([]T, len(arreglo))
		copy(heapArr.datos, arreglo)
	}

	heapArr.cantidad = len(arreglo)
	heapArr.heapify(heapArr.datos, funcion_cmp)
	return heapArr
}

// HEAPSORT -----------------------------------------------------------------------------------------------------------

func (heap *colaPrioridad[T]) ordenarArreglo() {
	if heap.cantidad == 0 {
		return
	}
	heap.cantidad--
	heap.swapIndices(0, heap.cantidad)
	heap.downheap(heap.datos[0], 0)
	heap.ordenarArreglo()
}

func HeapSort[T comparable](elementos []T, funcion_cmp func(T, T) int) []T {
	heapArr := new(colaPrioridad[T])
	heapArr.cantidad = len(elementos)
	heapArr.datos = elementos
	heapArr.heapify(heapArr.datos, funcion_cmp)
	heapArr.ordenarArreglo()
	return elementos
}

// AUXILIARES ---------------------------------------------------------------------------------------------------------

func (heap *colaPrioridad[T]) redimensionar(nuevoTamanio int) {
	temp := heap.datos
	heap.datos = make([]T, nuevoTamanio)
	copy(heap.datos, temp)
}

// Calculo de padres e hijos >>>>>>>>>>>>>>>>>>>>>>>>>>>

// hijoMax recibe dos elementos en un arreglo y sus posiciones y devuelve la posicion del mayor.
// En caso de que sean iguales, devuelve la posición del primero.
func (heap colaPrioridad[T]) hijoMax(x T, y T, posX int, posY int) int {
	if heap.cmp(y, x) > 0 && posY < heap.cantidad {
		return posY
	}
	return posX
}

// abs devuelve el valor absoluto de un integer
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func padre(indice int) int {
	return abs((indice - 1) / 2)
}

func hijoIzq(indice int) int {
	return 2*indice + 1
}

func hijoDer(indice int) int {
	return 2*indice + 2
}

// swapIndices recibe dos indices de un array e intercambia los valores guardados en el indice
func (heap *colaPrioridad[T]) swapIndices(i, j int) {
	heap.datos[i], heap.datos[j] = heap.datos[j], heap.datos[i]
}

// condicionHeap devuelve True si el padre es mayor o igual al hijo que estamos evaluando actualmente
func (heap colaPrioridad[T]) condicionHeap(actual, padre T) bool {
	return heap.cmp(actual, padre) < 0 || actual == padre
}

func (heap *colaPrioridad[T]) upheap(elem T, pos int) {
	if pos == 0 {
		return
	}

	if !heap.condicionHeap(elem, heap.datos[padre(pos)]) {
		heap.swapIndices(pos, padre(pos))
		heap.upheap(elem, padre(pos))
	}
}

func (heap *colaPrioridad[T]) downheap(elem T, pos int) {
	if pos == heap.cantidad-1 {
		return
	}

	var posHijoIzq = hijoIzq(pos)
	var posHijoDer = hijoDer(pos)

	// es necesario controlar que las posiciones calculadas para los hijos esten dentro de los lim. del arreglo
	if posHijoIzq >= heap.cantidad {
		return
	}

	nuevoIndice := posHijoIzq
	if posHijoDer < heap.cantidad {
		nuevoIndice = heap.hijoMax(heap.datos[posHijoIzq], heap.datos[posHijoDer], posHijoIzq, posHijoDer)
	}

	if !heap.condicionHeap(heap.datos[nuevoIndice], elem) {
		heap.swapIndices(pos, nuevoIndice)
		heap.downheap(elem, nuevoIndice)
	}

}

// PRIMITIVAS COLA DE PRIORIDAD -------------------------------------------------------------------------------------

// EstaVacia devuelve true si la cola se encuentra vacía, false en caso contrario.
func (heap colaPrioridad[T]) EstaVacia() bool {
	return heap.cantidad == 0
}

// Encolar Agrega un elemento al heap.
func (heap *colaPrioridad[T]) Encolar(elem T) {
	if len(heap.datos) == heap.cantidad {
		heap.redimensionar(heap.cantidad * _FACTOR_REDIMENSION)
	}

	heap.datos[heap.cantidad] = elem
	heap.upheap(elem, heap.cantidad)
	heap.cantidad++
}

// VerMax devuelve el elemento con máxima prioridad. Si está vacía, entra en pánico con un mensaje
// "La cola esta vacia".
func (heap colaPrioridad[T]) VerMax() T {
	if heap.EstaVacia() {
		panic("La cola esta vacia")
	}
	return heap.datos[0]
}

// Desencolar elimina el elemento con máxima prioridad, y lo devuelve. Si está vacía, entra en pánico con un
// mensaje "La cola esta vacia"
func (heap *colaPrioridad[T]) Desencolar() T {
	if heap.EstaVacia() {
		panic("La cola esta vacia")
	}

	heap.cantidad--
	heap.swapIndices(0, heap.cantidad)
	heap.downheap(heap.datos[0], 0)

	if (len(heap.datos) >= heap.cantidad*_PROPORCION_MIN) && (len(heap.datos)/_FACTOR_REDIMENSION >= _TAMANIO_INICIAL) {
		heap.redimensionar(len(heap.datos) / _FACTOR_REDIMENSION)
	}

	return heap.datos[heap.cantidad]
}

// Cantidad devuelve la cantidad de elementos que hay en la cola de prioridad.
func (heap colaPrioridad[T]) Cantidad() int {
	return heap.cantidad
}
