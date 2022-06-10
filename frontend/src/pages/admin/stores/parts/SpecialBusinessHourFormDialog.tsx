import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import SpecialBusinessHourForm from "./SpecialBusinessHourForm";
import { SpecialBusinessHour } from "../../../../libs/SpecialBusinessHour";
import { BusinessHour } from "../../../../libs/BusinessHour";

type SpecialBusinessHourFormDialogProps = {
  editItem: SpecialBusinessHour;
  hours: BusinessHour[];
  open: boolean;
  onClose: (item: SpecialBusinessHour | null) => void;
};

export default function SpecialBusinessHourFormDialog(props: SpecialBusinessHourFormDialogProps) {
  const onSubmit = (data: SpecialBusinessHour) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
        <SpecialBusinessHourForm
            editItem={props.editItem}
            hours={props.hours}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></SpecialBusinessHourForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
