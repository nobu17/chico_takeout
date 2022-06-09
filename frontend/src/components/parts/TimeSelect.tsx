import { MenuItem, TextField, TextFieldProps } from "@mui/material";
import type { ChangeEventHandler, FocusEventHandler } from "react";

export type TimeSelectProps = {
  label: string;
  timeList: string[];
  error?: string;
};

export function TimeSelect(
  props: TimeSelectProps & {
    inputRef: TextFieldProps["ref"];
    value: string; // HH:mm
    onChange: ChangeEventHandler<HTMLTextAreaElement>;
    onBlur: FocusEventHandler<HTMLTextAreaElement>;
  }
) {
  return (
    <>
      <TextField
        ref={props.inputRef}
        value={props.value}
        onChange={props.onChange}
        onBlur={props.onBlur}
        select
        sx={{ mt: 2 }}
        fullWidth
        label={props.label}
        error={Boolean(props.error)}
        helperText={props.error}
      >
        {props.timeList.map((time, index) => (
          <MenuItem key={index} value={time}>
            {time}
          </MenuItem>
        ))}
      </TextField>
    </>
  );
}
