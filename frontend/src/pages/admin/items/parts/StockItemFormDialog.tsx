import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import StockItemForm from "./StockItemForm";
import { ItemKind } from "../../../../libs/ItemKind";
import { StockItem } from "../../../../libs/StockItem";

type StockItemFormDialogProps = {
  editItem: StockItem;
  itemKinds: ItemKind[],
  open: boolean;
  onClose: (item: StockItem | null) => void;
};

export default function StockItemFormDialog(props: StockItemFormDialogProps) {
  const onSubmit = (data: StockItem) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <StockItemForm
            itemKinds={props.itemKinds}
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></StockItemForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
