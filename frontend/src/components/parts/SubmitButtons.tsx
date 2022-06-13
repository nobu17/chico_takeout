import { Button } from "@mui/material";

type SubmitButtonsProps = {
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (): void;
}
interface callbackCancel {
  (): void;
}

export default function SubmitButtons(props: SubmitButtonsProps) {
  return (
    <>
      <Button
        color="primary"
        variant="contained"
        size="large"
        onClick={props.onSubmit}
      >
        確定
      </Button>
      <Button
        color="error"
        variant="contained"
        size="large"
        onClick={props.onCancel}
      >
        キャンセル
      </Button>
    </>
  );
}
