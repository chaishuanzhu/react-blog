declare module 'jinrishici' {
  export function load(
    onSuccess: (res: { data: { content: string } }) => void,
    onError?: (err: unknown) => void
  ): void;
}
