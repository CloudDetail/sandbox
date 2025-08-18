class User {
    constructor(id, name, email) {
        this.id = id;
        this.name = name;
        this.email = email;
    }

    static fromJSON(json) {
        return new User(json.id, json.name, json.email);
    }

    toJSON() {
        return {
            id: this.id,
            name: this.name,
            email: this.email
        };
    }
}

module.exports = { User };
