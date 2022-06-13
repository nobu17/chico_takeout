import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { RhfTimeSelect } from "../../../../components/parts/Rhf/RhfTimeSelect";
import { SpecialBusinessHour } from "../../../../libs/SpecialBusinessHour";
import { BusinessHour } from "../../../../libs/BusinessHour";
import { StoreTimeList } from "../../../../libs/Constant";
import { RhfSelects } from "../../../../components/parts/Rhf/RhfSelects";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
} from "../../../../libs/ErrorMessages";

import { startIsLessThanEnd } from "../../../../libs/util/TimeCompare";
import { ToDate, ToDateString } from "../../../../libs/util/DateUtil";

type SpecialBusinessHourFormProps = {
  editItem: SpecialBusinessHour;
  hours: BusinessHour[];
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (item: SpecialBusinessHour): void;
}
interface callbackCancel {
  (): void;
}

type SpecialBusinessHourInput = {
  id: string;
  name: string;
  date: Date;
  start: string;
  end: string;
  businessHourId: string;
};

const convertInput = (item: SpecialBusinessHour): SpecialBusinessHourInput => {
  return {
    id: item.id,
    name: item.name,
    date: ToDate(item.date),
    start: item.start,
    end: item.end,
    businessHourId: item.businessHourId,
  };
};

const reverseInput = (item: SpecialBusinessHourInput): SpecialBusinessHour => {
  return {
    id: item.id,
    name: item.name,
    date: ToDateString(item.date),
    start: item.start,
    end: item.end,
    businessHourId: item.businessHourId,
  };
};

type HourSelects = {
  name: string;
  value: string;
};

const convertHoursList = (hours: BusinessHour[]): HourSelects[] => {
  return hours.map((val) => ({
    name: val.name,
    value: val.id,
  }));
};

export default function SpecialBusinessHourForm(
  props: SpecialBusinessHourFormProps
) {
  const {
    register,
    control,
    handleSubmit,
    setError,
    setValue,
    formState: { errors },
  } = useForm<SpecialBusinessHourInput>({
    defaultValues: convertInput(props.editItem),
  });
  const onSubmit: SubmitHandler<SpecialBusinessHourInput> = (data) => {
    if (!startIsLessThanEnd(data.start, data.end)) {
      setError(`start`, { message: "開始時刻が終了時刻よりも大きいです。" });
      setError(`end`, { message: "開始時刻が終了時刻よりも大きいです。" });
      return;
    }
    const converted = reverseInput(data);
    props.onSubmit(converted);
  };
  const onCancel = () => {
    props.onCancel();
  };
  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 10, message: MaxLengthErrorMessage(10) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <RhfSelects
            label="販売時間"
            name="businessHourId"
            multiple={false}
            itemList={convertHoursList(props.hours)}
            control={control}
          />
          <Controller
            name="date"
            control={control}
            render={({ field }) => {
              return (
                <MobileDatePicker
                  {...field}
                  label="営業時間"
                  mask="____/__/__ __:__:__"
                  inputFormat="yyyy/MM/dd"
                  renderInput={(props) => <TextField {...props} />}
                  onChange={(newValue) => {
                    if (newValue != null) {
                      setValue("date", newValue);
                    }
                  }}
                />
              );
            }}
          />
          <RhfTimeSelect
            label="開始時間"
            timeList={timeList}
            name="start"
            control={control}
          />
          <RhfTimeSelect
            label="終了時間"
            timeList={timeList}
            name="end"
            control={control}
          />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}

const timeList = StoreTimeList;
