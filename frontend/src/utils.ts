interface DebouncedFunction<T extends (...args: any[]) => any> {
    (...args: Parameters<T>): void;
    cancel: () => void;
}

function useDebounce<T extends (...args: any[]) => any>(func: T, delay: number): DebouncedFunction<T> {
    let timeoutId: ReturnType<typeof setTimeout> | null;

    const debounced = (...args: Parameters<T>) => {
        clearTimeout(timeoutId!);
        timeoutId = setTimeout(() => {
            func(...args);
        }, delay);
    };

    debounced.cancel = () => {
        clearTimeout(timeoutId!);
    };

    return debounced;
}

export default useDebounce;
