// Modern JavaScript Features Test

// Arrow functions and template literals
const greet = (name) => {
    const time = new Date().getHours();
    return `Good ${time < 12 ? 'morning' : 'evening'}, ${name}!`;
};

// Destructuring and spread operator
const config = {
    theme: 'dark',
    language: 'en',
    notifications: true
};

const { theme, ...rest } = config;
const fullConfig = {
    ...rest,
    theme: 'light',
    version: '2.0'
};

// Async/await and promises
async function fetchUserData(userId) {
    try {
        const response = await fetch(`/api/users/${userId}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error(`Error fetching user: ${error.message}`);
        return null;
    }
}

// Class with private fields and methods
class UserManager {
    #users = new Map();
    #lastId = 0;

    constructor(initialUsers = []) {
        initialUsers.forEach(user => this.#addUser(user));
    }

    #addUser(userData) {
        this.#lastId++;
        this.#users.set(this.#lastId, userData);
        return this.#lastId;
    }

    addNewUser(userData) {
        return this.#addUser(userData);
    }

    get userCount() {
        return this.#users.size;
    }
}
