import { MenuItem, TextField, TextFieldProps } from "@mui/material";
import type { ChangeEventHandler, FocusEventHandler } from "react";

import {
  HourOffset,
  HourOffsets
} from "../../libs/util/HourOffset";

export type HourOffsetSelectProps = {
  label: string;
  error?: string;
};

export function HourOffsetSelect(
  props: HourOffsetSelectProps & {
    inputRef: TextFieldProps["ref"];
    value: HourOffset[];
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
          multiple: false,
          value: props.value,
        }}
        sx={{ mt: 2 }}
        fullWidth
        label={props.label}
        error={Boolean(props.error)}
        helperText={props.error}
      >
        {HourOffsets.map((hour, index) => (
          <MenuItem key={index} value={hour}>
            {hour}
          </MenuItem>
        ))}
      </TextField>
    </>
  );
}
