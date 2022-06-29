import * as React from "react";
import { Stack, TextField, Button } from "@mui/material";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import { RhfSelects } from "../../../components/parts/Rhf/RhfSelects";

import { PickupDate } from "../../../hooks/UsePickupDate";
import {
  ToDateString,
  GetTimeList,
  IsFutureFromNow,
  GetDateTimeFromStr,
} from "../../../libs/util/DateUtil";

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
};

type SelectValue = {
  value: string;
};

export default function PickupSelect(props: PickupSelectProps) {
  const { control, handleSubmit, setError, setValue } = useForm<PickupDate>({
    defaultValues: props.selectedInfo,
  });

  const [timeList, setTimeList] = React.useState<SelectValue[]>([]);
  // for retain date picker value
  // const [selectedDate, setSelectedDate] = React.useState<Date>();
  const onSubmit: SubmitHandler<PickupDate> = (data) => {
    setValue("time", data.time);
    const inputDateTime = GetDateTimeFromStr(data.date, data.time);
    // check pick up time is over 1 hour later
    if (!IsFutureFromNow(inputDateTime, 60)) {
      setError("time", {
        message:
          "現在日時から近すぎるため、お手数ですが余裕を持った時間を指定してください。(お手数ですが画面をリロードすることをオススメします。)",
      });
      return;
    }
    if (props.onSubmit) {
      props.onSubmit(data);
    }
  };

  React.useEffect(() => {
    updateTimeList(props.selectedInfo.date);
    setValue("time", props.selectedInfo.time);
  }, []);

  const handleSelectDate = (date: Date) => {
    const dateStr = ToDateString(date);
    setValue("date", dateStr);
    // reset selectable times
    setValue("time", "");
    updateTimeList(dateStr);
  };

  const updateTimeList = (dateStr: string) => {
    const sameDates = props.selectableInfo
      .filter((item) => item.date === dateStr)
      .sort((a, b) => {
        if (a.startTime < b.startTime) return -1;
        if (a.startTime > b.startTime) return 1;
        return 0;
      });
    if (sameDates) {
      // list up all times
      let times: string[] = [];
      for (const date of sameDates) {
        const time = GetTimeList(date.startTime, date.endTime, 30);
        times = times.concat(time);
      }
      setTimeList(times.map((t) => ({ value: t })));
    } else {
      setTimeList([]);
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
              label="受取日時"
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
                const dateStr = ToDateString(day);
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
      <RhfSelects
        label="受取時間帯"
        name="time"
        multiple={false}
        itemList={timeList}
        control={control}
      />
      <Stack direction="row" spacing={2}>
        <Button onClick={handleSubmit(onSubmit)} variant="contained">
          次へ
        </Button>
      </Stack>
    </Stack>
  );
}
