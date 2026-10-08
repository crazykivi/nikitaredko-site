import type { ComputedRef, Ref } from "vue";
import type { Command } from "../types/commandPalette";

interface KeyboardOptions {
    isOpen: Ref<boolean>;
    selectedIndex: Ref<number>;
    flatCommands: ComputedRef<Command[]>;
    onExecute: (cmd: Command) => void;
    onOpen: () => void;
    onClose: () => void;
    onQuickNavigate: (index: number) => void;
}

function isTypingField(target: EventTarget | null): boolean {
    if (!(target instanceof HTMLElement)) return false;

    const tag = target.tagName.toLowerCase();

    return (
        tag === "input" ||
        tag === "textarea" ||
        tag === "select" ||
        target.isContentEditable
    );
}

export function useCommandPaletteKeyboard(options: KeyboardOptions) {
    const {
        isOpen,
        selectedIndex,
        flatCommands,
        onExecute,
        onOpen,
        onClose,
        onQuickNavigate,
    } = options;

    const moveSelection = (delta: number) => {
        const len = flatCommands.value.length;
        if (len === 0) return;

        selectedIndex.value = (selectedIndex.value + delta + len) % len;
    };

    const selectIndex = (index: number) => {
        const len = flatCommands.value.length;
        if (len === 0) return;

        selectedIndex.value = Math.min(Math.max(index, 0), len - 1);
    };

    const executeSelected = () => {
        const cmd = flatCommands.value[selectedIndex.value];
        if (cmd) onExecute(cmd);
    };

    const handleGlobalKeydown = (e: KeyboardEvent) => {
        const isMac = navigator.platform.toUpperCase().includes("MAC");
        const mod = isMac ? e.metaKey : e.ctrlKey;
        const typing = isTypingField(e.target);

        if (mod && e.shiftKey && e.key.toLowerCase() === "k") {
            e.preventDefault();
            isOpen.value ? onClose() : onOpen();
            return;
        }

        if (!isOpen.value) {
            if (e.key === "/" && !mod && !e.altKey && !e.shiftKey) {
                if (typing) return;

                e.preventDefault();
                onOpen();
                return;
            }

            if (mod && !e.shiftKey && /^[1-4]$/.test(e.key)) {
                if (typing) return;

                e.preventDefault();
                onQuickNavigate(Number(e.key) - 1);
                return;
            }

            return;
        }

        if (e.key === "Escape") {
            e.preventDefault();
            onClose();
            return;
        }

        switch (e.key) {
            case "ArrowDown":
                e.preventDefault();
                moveSelection(1);
                break;

            case "ArrowUp":
                e.preventDefault();
                moveSelection(-1);
                break;

            case "Tab":
                e.preventDefault();
                moveSelection(e.shiftKey ? -1 : 1);
                break;

            case "Home":
                e.preventDefault();
                selectIndex(0);
                break;

            case "End":
                e.preventDefault();
                selectIndex(flatCommands.value.length - 1);
                break;

            case "Enter":
                e.preventDefault();
                executeSelected();
                break;
        }
    };

    return {
        handleGlobalKeydown,
    };
}