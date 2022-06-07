import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import StockItemRemainForm from "./StockItemRemainForm";

type StockItemRemainFormDialogProps = {
  editItem: StockItemRemain;
  open: boolean;
  onClose: (item: StockItemRemain | null) => void;
};

type StockItemRemain = {
  id: string;
  name: string;
  remain: number;
};

export default function StockItemRemainFormDialog(
  props: StockItemRemainFormDialogProps
) {
  const onSubmit = (data: StockItemRemain) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <StockItemRemainForm
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></StockItemRemainForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
