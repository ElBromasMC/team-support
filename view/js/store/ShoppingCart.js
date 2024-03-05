
export default class ShoppingCart extends HTMLElement {
    constructor() {
        super();
        this.CLIENTS = [];
        this.items = JSON.parse(localStorage.getItem("cartItems")) || [];
    }
    static get observedAttributes() {
        return [];
    }
    connectedCallback() {
        this.addEventListener("NEWCLIENT", e => {
            // e.detail is an object that implements EventTarget and UPDATE event
            const client = e.detail;
            this.CLIENTS.push(client);
            client.dispatchEvent(new CustomEvent("UPDATE", { detail: _ => structuredClone(this.items) }));
        });
        this.addEventListener("UPDATE", e => {
            // e.detail is a function that updates cart items
            const newItems = e.detail(structuredClone(this.items));
            console.log(newItems);
            this.items = newItems;
            localStorage.setItem("cartItems", JSON.stringify(newItems));
            for (let client of this.CLIENTS) {
                client.dispatchEvent(new CustomEvent("UPDATE", { detail: _ => structuredClone(newItems) }));
            }
        });
    }
}
