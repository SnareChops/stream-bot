export type FlatObject = { [key: string]: string | number | boolean };

export function createElement<T extends HTMLElement>(tag: string, attrs: FlatObject = {}): T {
    const element = document.createElement(tag);
    for (const key of Object.keys(attrs)) {
        element.setAttribute(key, attrs[key].toString());
    }
    return element as T;
}

export function createComponent<T extends HTMLElement>(component: { new(): T, get tag(): string }, attrs: FlatObject = {}): T {
    return createElement(component.tag, attrs)
}

export function Component(ctr: Function & { get tag(): string }, _: ClassDecoratorContext) {
    window.customElements.define(ctr.tag, ctr as unknown as typeof HTMLElement);
}
