<template>
  <section>
    <div
      v-for="(todo, t) of todos"
      :key="t"
      :class="{ 'completed': todo.completed }"
    >
      <form @submit.prevent="update(todo)">
        <input
          type="checkbox"
          :value="todo.completed"
          :checked="todo.completed"
          @input="check(todo)"
        >
        <input
          type="text"
          class="editable"
          v-model="todo.content"
          @keyup.enter="$event.target.blur()"
        >
      </form>
    </div>
    <br>
    <div>
      <form @submit.prevent="add">
        <input
          type="text"
          v-model="newTodo.content"
        >
        &nbsp;
        <button type="submit">Add</button>
      </form>
    </div>
  </section>
</template>

<script>
export default {
  data: () => ({
    todos: [],
    newTodo: {
      content: '',
    },
  }),

  created() {
    this.all();
  },

  methods: {
    all() {
      fetch('http://localhost:4321/todos')
        .then(res => res.json())
        .then(data => {
          this.todos = data;
        });
    },
    add() {
      fetch('http://localhost:4321/todos', {
        method: 'POST',
        body: JSON.stringify(this.newTodo),
      }).then(() => {
        this.all();
        this.newTodo = this.$options.data().newTodo;
      });
    },
    update(todo) {
      fetch(`http://localhost:4321/todos/${todo.id}`, {
        method: 'PUT',
        body: JSON.stringify(todo),
      }).then(() => {
        this.all();
      });
    },
    check(todo) {
      todo.completed = !todo.completed;
      this.update(todo);
    },
  },
};
</script>

<style>
section { display: flex; width: 400px; margin: auto; flex-direction: column; align-items: center; }
.editable { border: transparent; }
.editable:focus { border: black; }
.completed .editable { text-decoration: line-through; }
</style>
