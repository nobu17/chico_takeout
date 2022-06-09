import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import BusinessHourForm from "./BusinessHourForm";
import { BusinessHour } from "../../../../libs/BusinessHour";

type BusinessHourFormDialogProps = {
  editItem: BusinessHour;
  open: boolean;
  onClose: (item: BusinessHour | null) => void;
};

export default function FoodItemFormDialog(props: BusinessHourFormDialogProps) {
  const onSubmit = (data: BusinessHour) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <BusinessHourForm
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></BusinessHourForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
