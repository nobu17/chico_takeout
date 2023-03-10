import * as React from "react";
import { Stack, TextField, Button, Typography } from "@mui/material";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import { RhfGroupSelects } from "../../../components/parts/Rhf/RhfGroupSelects";

import { PickupDate } from "../../../hooks/UsePickupDate";
import {
  toDateString,
  GetTimeList,
  isBeforeFromNow,
  GetDateTimeFromStr,
} from "../../../libs/util/DateUtil";
import { OffsetMinutesUserCanOrder, OffsetMinutesUserCanCancel } from "../../../libs/Constant";

type PickupSelectProps = {
  selectedInfo: PickupDate;
  selectableInfo: SelectableDateTimeInfo[];
  onSubmit?: callbackSubmit;
};
interface callbackSubmit {
  (selectedInfo: PickupDate): void;
}

type SelectableDateTimeInfo = {
  date: string; // yyyy/MM/dd
  startTime: string; // HH:MM
  endTime: string; // HH:MM
  hourName: string;
};

type SelectTimeValue = {
  groupName: string;
  value: string[];
};

export default function PickupSelect(props: PickupSelectProps) {
  const { control, handleSubmit, setError, setValue } = useForm<PickupDate>({
    defaultValues: props.selectedInfo,
  });

  const [categoryTimes, setCategoryTimes] = React.useState<SelectTimeValue[]>(
    []
  );
  // for retain date picker value
  // const [selectedDate, setSelectedDate] = React.useState<Date>();
  const onSubmit: SubmitHandler<PickupDate> = (data) => {
    setValue("time", data.time);
    const inputDateTime = GetDateTimeFromStr(data.date, data.time);
    // check pick up time is over 1 hour from now
    if (isBeforeFromNow(inputDateTime, OffsetMinutesUserCanOrder)) {
      setError("time", {
        message:
          "現在日時から近すぎるため、お手数ですが余裕を持った時間を指定してください。(お手数ですが画面をリロードをお願いいたします。)",
      });
      return;
    }
    if (props.onSubmit) {
      props.onSubmit(data);
    }
  };

  React.useEffect(() => {
    updateCategoryTimes(props.selectedInfo.date);
    setValue("time", props.selectedInfo.time);
  }, []);

  const handleSelectDate = (date: Date) => {
    const dateStr = toDateString(date);
    setValue("date", dateStr);
    // reset selectable times
    setValue("time", "");
    updateCategoryTimes(dateStr);
  };

  const updateCategoryTimes = (dateStr: string) => {
    const sameDates = props.selectableInfo
      .filter((item) => item.date === dateStr)
      .sort((a, b) => {
        if (a.startTime < b.startTime) return -1;
        if (a.startTime > b.startTime) return 1;
        return 0;
      });
    if (sameDates) {
      // list up all times
      let times: SelectTimeValue[] = [];
      for (const date of sameDates) {
        const timeList = GetTimeList(date.startTime, date.endTime, 30);
        times.push({ groupName: date.hourName, value: timeList });
      }
      setCategoryTimes(times);
    } else {
      setCategoryTimes([]);
    }
  };

  return (
    <Stack spacing={3}>
      <Controller
        name="date"
        control={control}
        rules={{ required: true }}
        render={({ field }) => {
          return (
            <MobileDatePicker
              {...field}
              label="来店日時"
              mask="____/__/__"
              inputFormat="yyyy/MM/dd"
              closeOnSelect={true}
              renderInput={(props) => <TextField {...props} />}
              onChange={(newValue) => {
                // this timing, not apply input. only retain value when close
                if (newValue != null) {
                  handleSelectDate(newValue);
                  // setSelectedDate(newValue);
                  // react hook form is needed to apply value selection
                  // const dateStr = ToDateString(newValue);
                  // setValue("date", dateStr);
                }
              }}
              onAccept={(val) => {
                // on close the dialog, apply value
                // if (selectedDate) {
                // handleSelectDate(selectedDate);
                // }
              }}
              disablePast={true}
              shouldDisableDate={(day: Date) => {
                // find from selectable date list
                const dateStr = toDateString(day);
                const date = props.selectableInfo.find(
                  (item) => item.date === dateStr
                );
                if (date) {
                  return false;
                }
                return true;
              }}
            />
          );
        }}
      />
      <RhfGroupSelects
        label="来店時間帯"
        name="time"
        multiple={false}
        itemList={categoryTimes}
        control={control}
      />
      <Typography variant="subtitle1" color="error" textAlign="center">
        ※注文後のキャンセルは、受け取り時間の{OffsetMinutesUserCanCancel / 60}時間前まで可能です。ご注意ください。
      </Typography>
      <Stack direction="row" spacing={2}>
        <Button
          onClick={handleSubmit(onSubmit)}
          variant="contained"
          sx={{ width: 100 }}
        >
          次へ
        </Button>
      </Stack>
    </Stack>
  );
}
