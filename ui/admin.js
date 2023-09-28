"use strict";
(() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __decorateClass = (decorators, target, key, kind) => {
    var result = kind > 1 ? void 0 : kind ? __getOwnPropDesc(target, key) : target;
    for (var i = decorators.length - 1, decorator; i >= 0; i--)
      if (decorator = decorators[i])
        result = (kind ? decorator(target, key, result) : decorator(result)) || result;
    if (kind && result)
      __defProp(target, key, result);
    return result;
  };
  var __accessCheck = (obj, member, msg) => {
    if (!member.has(obj))
      throw TypeError("Cannot " + msg);
  };
  var __privateGet = (obj, member, getter) => {
    __accessCheck(obj, member, "read from private field");
    return getter ? getter.call(obj) : member.get(obj);
  };
  var __privateAdd = (obj, member, value) => {
    if (member.has(obj))
      throw TypeError("Cannot add the same private member more than once");
    member instanceof WeakSet ? member.add(obj) : member.set(obj, value);
  };
  var __privateSet = (obj, member, value, setter) => {
    __accessCheck(obj, member, "write to private field");
    setter ? setter.call(obj, value) : member.set(obj, value);
    return value;
  };
  var __privateMethod = (obj, member, method) => {
    __accessCheck(obj, member, "access private method");
    return method;
  };

  // ui/src/lib/index.ts
  function createElement(tag, attrs = {}) {
    const element = document.createElement(tag);
    for (const key of Object.keys(attrs)) {
      element.setAttribute(key, attrs[key].toString());
    }
    return element;
  }
  function createComponent(component, attrs = {}) {
    return createElement(component.tag, attrs);
  }
  function Component(ctr, _) {
    window.customElements.define(ctr.tag, ctr);
  }

  // ui/src/admin/chatters-panel.ts
  var _list;
  var ChattersPanel = class extends HTMLElement {
    constructor() {
      super(...arguments);
      __privateAdd(this, _list, void 0);
    }
    static get tag() {
      return "chatters-panel";
    }
    connectedCallback() {
      __privateSet(this, _list, createElement("ul"));
      this.append(__privateGet(this, _list));
    }
  };
  _list = new WeakMap();
  ChattersPanel = __decorateClass([
    Component
  ], ChattersPanel);

  // ui/src/admin/quick-panel.ts
  var QuickPanel = class extends HTMLElement {
    static get tag() {
      return "quick-panel";
    }
    connectedCallback() {
      const button1 = createElement("button");
      button1.innerText = "Button 1";
      const button2 = createElement("button");
      button2.innerText = "Button 2";
      const button3 = createElement("button");
      button3.innerText = "Button 3";
      const button4 = createElement("button");
      button4.innerText = "Button 4";
      this.append(
        button1,
        button2,
        button3,
        button4
      );
    }
  };
  QuickPanel = __decorateClass([
    Component
  ], QuickPanel);

  // ui/src/admin/mystery-panel.ts
  var _socket, _onSignal, onSignal_fn, _loadItems, loadItems_fn;
  var MysteryPanel = class extends HTMLElement {
    constructor() {
      super(...arguments);
      __privateAdd(this, _onSignal);
      __privateAdd(this, _loadItems);
      __privateAdd(this, _socket, void 0);
    }
    static get tag() {
      return "mystery-panel";
    }
    static create(socket) {
      const el = createComponent(MysteryPanel);
      __privateSet(el, _socket, socket);
      return el;
    }
    connectedCallback() {
      __privateGet(this, _socket).addHandler(__privateMethod(this, _onSignal, onSignal_fn).bind(this));
      __privateGet(this, _socket).send({ kind: "admin.mystery.request" });
    }
    select(name) {
      __privateGet(this, _socket).send({ kind: "admin.mystery.select", args: [name] });
    }
  };
  _socket = new WeakMap();
  _onSignal = new WeakSet();
  onSignal_fn = function(signal) {
    switch (signal.kind) {
      case "admin.mystery.load":
        return __privateMethod(this, _loadItems, loadItems_fn).call(this, signal);
      default:
        return;
    }
  };
  _loadItems = new WeakSet();
  loadItems_fn = function(signal) {
    for (const child of this.children) {
      child.remove();
    }
    for (const item of signal.items) {
      this.append(MysteryListItem.create(this, item));
    }
  };
  MysteryPanel = __decorateClass([
    Component
  ], MysteryPanel);
  var _panel, _checkbox, _item, _changed, changed_fn;
  var MysteryListItem = class extends HTMLElement {
    constructor() {
      super(...arguments);
      __privateAdd(this, _changed);
      __privateAdd(this, _panel, void 0);
      __privateAdd(this, _checkbox, void 0);
      __privateAdd(this, _item, void 0);
    }
    static get tag() {
      return "mystery-list-item";
    }
    static create(panel, item) {
      const el = createComponent(MysteryListItem);
      __privateSet(el, _panel, panel);
      __privateSet(el, _item, item);
      return el;
    }
    connectedCallback() {
      __privateSet(this, _checkbox, createElement("input", { type: "checkbox" }));
      __privateGet(this, _checkbox).checked = __privateGet(this, _item).selected;
      __privateGet(this, _checkbox).onchange = __privateMethod(this, _changed, changed_fn).bind(this);
      const nameSpan = createElement("span");
      nameSpan.innerText = __privateGet(this, _item).name;
      const rangeSpan = createElement("span");
      rangeSpan.innerText = __privateGet(this, _item).start + "-" + __privateGet(this, _item).end;
      const speedSpan = createElement("span");
      speedSpan.innerText = "" + __privateGet(this, _item).speed;
      this.append(
        __privateGet(this, _checkbox),
        nameSpan,
        rangeSpan,
        speedSpan
      );
    }
  };
  _panel = new WeakMap();
  _checkbox = new WeakMap();
  _item = new WeakMap();
  _changed = new WeakSet();
  changed_fn = function() {
    if (__privateGet(this, _checkbox).checked !== __privateGet(this, _item).selected) {
      __privateGet(this, _item).selected = __privateGet(this, _checkbox).checked;
      __privateGet(this, _panel).select(__privateGet(this, _item).name);
    }
  };
  MysteryListItem = __decorateClass([
    Component
  ], MysteryListItem);

  // ui/src/lib/socket.ts
  var Socket = class {
    #url;
    #socket;
    #handlers = [];
    #interval = 1e3;
    #maxInterval = 3e4;
    #decay = 2;
    constructor(url) {
      this.#url = url;
      this.#connect();
    }
    addHandler(handler) {
      this.#handlers.push(handler);
    }
    send(signal) {
      this.#socket?.send(JSON.stringify(signal));
    }
    #connect() {
      this.#socket = new WebSocket(this.#url);
      this.#socket.onopen = () => {
        console.log("Socket connection open");
        this.#interval = 1e3;
      };
      this.#socket.onclose = (event) => {
        console.log("Socket connection closed", event);
        this.#reconnect();
      };
      this.#socket.onerror = (err) => console.error("Socket connection error", err);
      this.#socket.onmessage = (event) => {
        console.log("Socket message received", event);
        try {
          const message = JSON.parse(event.data);
          for (const handler of this.#handlers) {
            handler(message);
          }
        } catch (err) {
          console.log("Failed to parse message data", err);
        }
      };
    }
    #reconnect() {
      setTimeout(() => {
        console.log(`Socket reconnecting... (delay: ${this.#interval}ms)`);
        this.#connect();
        this.#interval = Math.min(
          this.#interval * this.#decay,
          this.#maxInterval
        );
      }, this.#interval);
    }
  };

  // ui/src/admin.ts
  var socketUrl = `ws://${location.host}/adminws`;
  var _socket2;
  var AdminPanel = class extends HTMLElement {
    constructor() {
      super(...arguments);
      __privateAdd(this, _socket2, void 0);
    }
    static get tag() {
      return "admin-panel";
    }
    connectedCallback() {
      __privateSet(this, _socket2, new Socket(socketUrl));
      this.append(
        createComponent(QuickPanel),
        createComponent(ChattersPanel),
        MysteryPanel.create(__privateGet(this, _socket2))
      );
    }
  };
  _socket2 = new WeakMap();
  AdminPanel = __decorateClass([
    Component
  ], AdminPanel);
  document.querySelector("body")?.appendChild(createComponent(AdminPanel));
})();
