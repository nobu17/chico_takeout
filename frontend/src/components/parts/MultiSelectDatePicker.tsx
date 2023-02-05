import * as React from "react";
import {
  Typography,
  IconButton,
  Stack,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import SubmitButtons from "./SubmitButtons";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import { Calendar, DateObject } from "react-multi-date-picker";
import DatePanel from "react-multi-date-picker/plugins/date_panel";

import { toDate, toDateString } from "../../libs/util/DateUtil";

type MultiSelectDatePickerProps = {
  selectedDates: string[];
  onSubmit: (selectedDates: string[]) => void;
  onCancel: () => void;
};

const convertStrToDate = (dates: string[]): Date[] => {
  return dates.map((d) => toDate(d));
};

const convertDateToStr = (dates: Date[]): string[] => {
  return dates.map((d) => toDateString(d));
};

const convertLabel = (selectedDates: string[]): string => {
  if (selectedDates.length <= 0) {
    return "指定なし";
  }

  let str = selectedDates.join(",");
  if (str.length > 15) {
    str = str.substr(0, 15) + "...";
  }
  return "(" + selectedDates.length + "個):" + str;
};

const months = [
  "1月",
  "2月",
  "3月",
  "4月",
  "5月",
  "6月",
  "7月",
  "8月",
  "9月",
  "10月",
  "11月",
  "12月",
];
const weekDays = ["日", "月", "火", "水", "木", "金", "土"];

export default function MultiSelectDatePicker(
  props: MultiSelectDatePickerProps
) {
  const [open, setOpen] = React.useState(false);
  const [values, setValues] = React.useState<Date[]>([]);

  const handleOpenEdit = () => {
    setValues(convertStrToDate(props.selectedDates));
    setOpen(true);
  };
  const handleSubmit = () => {
    const dates = convertDateToStr(values);
    setOpen(false);
    props.onSubmit(dates);
  };
  const handleCancel = () => {
    setOpen(false);
    props.onCancel();
  };

  const onCalendarChanged = (dates: DateObject[]) => {
    setValues(dates.map((d) => d.toDate()));
  };

  return (
    <>
      <Typography color="text.primary" sx={{ mb: 0 }}>
        販売期間指定
      </Typography>
      <Stack direction="row">
        <Typography variant="body2" sx={{ m: 1 }}>
          {convertLabel(props.selectedDates)}
        </Typography>
        <IconButton color="secondary" onClick={handleOpenEdit}>
          <CalendarMonthIcon />
        </IconButton>
      </Stack>
      <Dialog open={open} fullScreen>
        <DialogTitle>販売期間指定</DialogTitle>
        <DialogContent>
          <Calendar
            multiple
            value={values}
            onChange={onCalendarChanged}
            plugins={[<DatePanel sort="date" position="right" />]}
            months={months}
            weekDays={weekDays}
          ></Calendar>
          <Button
            sx={{ m: 2 }}
            variant="outlined"
            color="error"
            onClick={() => setValues([])}
          >
            空にする
          </Button>
        </DialogContent>
        <DialogActions>
          <SubmitButtons
            onCancel={handleCancel}
            onSubmit={handleSubmit}
          ></SubmitButtons>
        </DialogActions>
      </Dialog>
    </>
  );
}
