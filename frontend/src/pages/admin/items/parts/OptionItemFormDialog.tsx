import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import OptionItemForm from "./OptionItemForm";
import { OptionItem } from "../../../../libs/OptionItem";

type OptionItemFormDialogProps = {
  editItem: OptionItem;
  open: boolean;
  onClose: (item: OptionItem | null) => void;
};

export default function OptionItemFormDialog(props: OptionItemFormDialogProps) {
  const onSubmit = (data: OptionItem) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <OptionItemForm
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></OptionItemForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
