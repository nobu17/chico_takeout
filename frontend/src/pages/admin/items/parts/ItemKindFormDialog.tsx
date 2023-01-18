import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import ItemKindForm from "./ItemKindForm";
import { OptionItem } from "../../../../libs/OptionItem";
import { ItemKind } from "../../../../libs/ItemKind";

type ItemKindFormDialogProps = {
  editItem: ItemKind;
  optionItems: OptionItem[];
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
        <DialogContent>
          <ItemKindForm
            editItem={props.editItem}
            optionItems={props.optionItems}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></ItemKindForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
