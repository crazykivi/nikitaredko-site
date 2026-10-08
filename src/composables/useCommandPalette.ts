import { ref } from "vue";

const HOLD_DELAY_MS = 400;

export function useCommandPalette() {
    const isOpen = ref(false);
    const query = ref("");
    const selectedIndex = ref(0);
    const expandedId = ref<string | null>(null);

    let previouslyFocused: HTMLElement | null = null;
    let holdTimer: ReturnType<typeof setTimeout> | null = null;
    let suppressClick = false;

    const clearHoldTimer = () => {
        if (holdTimer !== null) {
            clearTimeout(holdTimer);
            holdTimer = null;
        }
    };

    const restoreFocus = () => {
        const el = previouslyFocused;
        previouslyFocused = null;

        if (el?.isConnected) {
            el.focus({ preventScroll: true });
            return;
        }

        if (document.activeElement instanceof HTMLElement) {
            document.activeElement.blur();
        }
    };

    const open = () => {
        previouslyFocused =
            document.activeElement instanceof HTMLElement && document.activeElement !== document.body
                ? document.activeElement
                : null;

        isOpen.value = true;
        query.value = "";
        selectedIndex.value = 0;
        expandedId.value = null;
        suppressClick = false;
        clearHoldTimer();
    };

    const close = () => {
        isOpen.value = false;
        expandedId.value = null;
        suppressClick = false;
        clearHoldTimer();
        restoreFocus();
    };

    const onCmdPointerDown = (id: string) => {
        suppressClick = false;
        clearHoldTimer();

        holdTimer = setTimeout(() => {
            suppressClick = true;
            expandedId.value = expandedId.value === id ? null : id;
        }, HOLD_DELAY_MS);
    };

    const onCmdPointerEnd = () => {
        clearHoldTimer();
    };

    const shouldSuppressClick = (): boolean => {
        if (suppressClick) {
            suppressClick = false;
            return true;
        }

        return false;
    };

    return {
        isOpen,
        query,
        selectedIndex,
        expandedId,
        open,
        close,
        onCmdPointerDown,
        onCmdPointerEnd,
        shouldSuppressClick,
    };
}