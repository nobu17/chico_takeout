import { useState } from "react";
import ConfirmDialog from "../components/parts/ConfirmDialog";

type State = {
  isOpen: boolean;
  title: string;
  message: string;
  resolve: (result: boolean) => void;
};

const initialState: State = {
  isOpen: false,
  title: "",
  message: "",
  resolve: () => {},
};
export function useConfirmDialog() {
  const [state, setState] = useState<State>(initialState);

  const showConfirmDialog = (title: string, message: string) => {
    const promise: Promise<boolean> = new Promise((resolve) => {
      const newState: State = {
        isOpen: true,
        title: title,
        message: message,
        resolve: resolve,
      };
      setState(newState);
    });
    return promise;
  };

  const handleClose = (result: boolean) => {
    state.resolve(result);
    setState(initialState);
  };

  const renderConfirmDialog = () => {
    return (
      <ConfirmDialog
        open={state.isOpen}
        title={state.title}
        message={state.message}
        onClose={handleClose}
      />
    );
  };

  return {
    showConfirmDialog,
    renderConfirmDialog,
  };
}
