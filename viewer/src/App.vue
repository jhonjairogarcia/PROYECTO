<template>
  <div class="container mx-auto p-4 h-screen">

    <!-- Encabezado de la página -->
    <header class="mb-4 bg-gray-300 py-6 border border-gray-400 rounded-md">
      <div class="mt-4 flex items-center">
        <h1 class="text-4xl font-bold text-blue-900 whitespace-nowrap mx-4">Search Email</h1>
        <input v-model="searchQuery" type="text" placeholder="Enter the text to search" @keyup.enter="searchPage" class="p-2 border border-gray-400 rounded-md text-lg w-full mx-4">
        <button @click="searchPage" :disabled="searching" class="px-4 py-2 bg-blue-500 text-white rounded-md text-lg mx-4">
          <span v-if="searching">Searching...</span>
          <span v-else>Search</span>
        </button>
      </div>
    </header>

    <!-- Botones de Paginación -->
    <div class="flex justify-center mb-4" v-if="searchResults.length > 0">
      <button @click="previousPage" :disabled="page === 0" class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md text-lg border border-gray-400">Previous</button>
      <span class="px-4 py-2 text-gray-700 rounded-md text-lg">Page: {{ page + 1 }}</span><!-- Mostrar el número de página actual -->
      <button @click="nextPage" class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md text-lg border border-gray-400">Next</button>
    </div>

    <!-- Resultados de la búsqueda -->
    <table v-if="searchResults.length" class="mb-4 w-full border-collapse border">
      <thead>
        <tr>
          <th class="px-4 py-2 border">Subject</th>
          <th class="px-4 py-2 border">From</th>
          <th class="px-4 py-2 border">To</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(email, index) in searchResults" :key="email._timestamp" @click="selectEmail(email)" :class="{'bg-gray-100': index % 2 === 0}">
          <td class="px-4 py-2 border">{{ email.Subject }}</td>
          <td class="px-4 py-2 border">{{ email.From }}</td>
          <td class="px-4 py-2 border">{{ email.To }}</td>
        </tr>
      </tbody>
    </table>

    <!-- Detalles del correo electrónico -->
    <div v-if="searchResults.length && selectedEmail" class="p-4 border border-gray-300 rounded-md">
      <h2 class="font-bold">Body</h2>
      {{ selectedEmail.Body }}
    </div>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        searchQuery: '',
        searchResults: [],
        selectedEmail: null,
        searching: false,
        end: false, // Para identificar si se llegó al final de la paginación
        page: 0 // Variable para la paginación
      };
    },
    methods: {
      async performSearch() {
        if (this.searchQuery.trim() === '') {
          // Si el campo de búsqueda está vacío, no se ejecuta la búsqueda
          return;
        }
        this.searching = true; // Activar el estado de búsqueda
        try {
          this.end = false; // Se inicia en falso para validar posteriormente si se alcanzó el final
          const response = await fetch(`http://localhost:8080/search?q=${this.searchQuery}&page=${this.page}`);
          if (response.ok) {
            const data = await response.json();
            if (data.length > 0) {
              this.searchResults = data; // Se obtienen los datos para visualizar en pantalla
            } else {
              if (this.page == 0){
                this.searchResults = data;
              }
              this.end = true; // Si se alcanzó el final
            }
            // Limpiar el contenido del correo electrónico seleccionado
            this.selectedEmail = null;
          } else {
            console.error('Error performing search:', response.statusText);
          }
        } catch (error) {
          console.error('Error performing search:', error);
        } finally {
          this.searching = false; // Desactivar el estado de búsqueda después de que finalice la búsqueda
        }
      },
      selectEmail(email) { // Selección de Email
        this.selectedEmail = email;
      },
      nextPage() { // Visualizar la página siguiente
        if (!this.end){
          this.page++;
          this.performSearch();
        }
      },
      previousPage() { // Visualizar la página anterior
        if (this.page > 0) {
          this.page--;
          this.performSearch();
        }
      },
      searchPage() { // Visualizar la primera página
        this.page = 0; // Reiniciar la página a 0
        this.performSearch();
      },
    }
  };
</script>

<style>
/* Estilos del componente App.vue */
</style>
