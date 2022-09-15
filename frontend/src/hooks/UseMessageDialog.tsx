import { useState } from "react";
import MessageDialog from "../components/parts/MessageDialog";

type State = {
  isOpen: boolean;
  title: string;
  message: string;
  resolve: () => void;
};

const initialState: State = {
  isOpen: false,
  title: "",
  message: "",
  resolve: () => {},
};
export function useMessageDialog() {
  const [state, setState] = useState<State>(initialState);

  const showMessageDialog = (title: string, message: string) => {
    const promise: Promise<void> = new Promise((resolve) => {
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

  const handleClose = () => {
    state.resolve();
    setState(initialState);
  };

  const renderDialog = () => {
    return (
      <MessageDialog
        open={state.isOpen}
        title={state.title}
        message={state.message}
        onClose={handleClose}
      />
    );
  };

  return {
    showMessageDialog,
    renderDialog,
  };
}
