/**
 * Todo List Application
 * This is a sample JavaScript file to test the minifier
 * It includes various JavaScript features like:
 * - Classes
 * - Arrow functions
 * - Template literals
 * - Local storage
 * - Event handling
 * - DOM manipulation
 */

class TodoApp {
    constructor() {
        this.todos = JSON.parse(localStorage.getItem('todos')) || [];
        this.todoInput = document.getElementById('todoInput');
        this.todoList = document.getElementById('todoList');
        this.statsElement = document.getElementById('stats');
        this.renderTodos();
        this.updateStats();
        
        // Add event listener for Enter key
        this.todoInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                this.addTodo();
            }
        });
    }

    // Add a new todo item
    addTodo() {
        const todoText = this.todoInput.value.trim();
        if (todoText) {
            const todo = {
                id: Date.now(),
                text: todoText,
                completed: false,
                createdAt: new Date().toISOString()
            };
            
            this.todos.push(todo);
            this.saveTodos();
            this.todoInput.value = '';
            this.renderTodos();
            this.updateStats();
        }
    }

    // Toggle todo completion status
    toggleTodo(id) {
        const todo = this.todos.find(todo => todo.id === id);
        if (todo) {
            todo.completed = !todo.completed;
            this.saveTodos();
            this.renderTodos();
            this.updateStats();
        }
    }

    // Delete a todo item
    deleteTodo(id) {
        this.todos = this.todos.filter(todo => todo.id !== id);
        this.saveTodos();
        this.renderTodos();
        this.updateStats();
    }

    // Save todos to localStorage
    saveTodos() {
        localStorage.setItem('todos', JSON.stringify(this.todos));
    }

    // Render the todo list
    renderTodos() {
        this.todoList.innerHTML = '';
        this.todos.forEach(todo => {
            const todoElement = document.createElement('div');
            todoElement.className = `todo-item ${todo.completed ? 'completed' : ''}`;
            
            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.checked = todo.completed;
            checkbox.onclick = () => this.toggleTodo(todo.id);
            
            const textSpan = document.createElement('span');
            textSpan.style.marginLeft = '10px';
            textSpan.textContent = todo.text;
            
            const deleteButton = document.createElement('button');
            deleteButton.textContent = 'Delete';
            deleteButton.onclick = () => this.deleteTodo(todo.id);
            
            todoElement.appendChild(checkbox);
            todoElement.appendChild(textSpan);
            todoElement.appendChild(deleteButton);
            
            this.todoList.appendChild(todoElement);
        });
    }

    // Update statistics
    updateStats() {
        const totalTodos = this.todos.length;
        const completedTodos = this.todos.filter(todo => todo.completed).length;
        const pendingTodos = totalTodos - completedTodos;
        
        this.statsElement.innerHTML = `
            <strong>Statistics:</strong><br>
            Total: ${totalTodos} | 
            Completed: ${completedTodos} | 
            Pending: ${pendingTodos}
        `;
    }
}

// Initialize the app
const todoApp = new TodoApp();

// Global function for the onclick handler
function addTodo() {
    todoApp.addTodo();
}
