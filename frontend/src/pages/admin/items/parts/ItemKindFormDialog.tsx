import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";

import ItemKindForm from "./ItemKindForm";
import { ItemKind } from "../../../../libs/ItemKind";

type ItemKindFormDialogProps = {
  editItem: ItemKind;
  open: boolean;
  onClose: (item: ItemKind | null) => void;
};

export default function ItemKindFormDialog(props: ItemKindFormDialogProps) {
  const onSubmit = (data: ItemKind) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogTitle>編集</DialogTitle>
        <DialogContent>
          <ItemKindForm
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></ItemKindForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
