import { MenuItem, TextField, TextFieldProps } from "@mui/material";
import type { ChangeEventHandler, FocusEventHandler } from "react";

import {
  DayOfWeek,
  toString,
  AllDayOfWeek,
} from "../../libs/util/DayOfWeek";

export type DayOfWeekSelectProps = {
  label: string;
  error?: string;
};

export function DayOfWeekSelect(
  props: DayOfWeekSelectProps & {
    inputRef: TextFieldProps["ref"];
    value: DayOfWeek[];
    onChange: ChangeEventHandler<HTMLTextAreaElement>;
    onBlur: FocusEventHandler<HTMLTextAreaElement>;
  }
) {
  return (
    <>
      <TextField
        ref={props.inputRef}
        onChange={props.onChange}
        onBlur={props.onBlur}
        select
        SelectProps={{
          multiple: true,
          value: props.value,
        }}
        sx={{ mt: 2 }}
        fullWidth
        label={props.label}
        error={Boolean(props.error)}
        helperText={props.error}
      >
        {AllDayOfWeek.map((dayOfWeek, index) => (
          <MenuItem key={index} value={dayOfWeek}>
            {toString(dayOfWeek)}
          </MenuItem>
        ))}
      </TextField>
    </>
  );
}
