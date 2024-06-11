export function Once<T extends (...args: any[]) => any>(fn: T): T {
    let called = false;
    let result: ReturnType<T>;

    return function (...args: Parameters<T>): ReturnType<T> {
        if (!called) {
            called = true;
            result = fn(...args);
        }
        return result;
    } as T;
}
