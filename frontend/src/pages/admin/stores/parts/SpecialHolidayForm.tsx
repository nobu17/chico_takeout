import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { SpecialHoliday } from "../../../../libs/SpecialHoliday";
import {
  SubmitHandler,
  useForm,
  FieldError,
  Controller,
} from "react-hook-form";

import { startDateIsLessThanEndDate } from "../../../../libs/util/TimeCompare";
import { ToDate, ToDateString } from "../../../../libs/util/DateUtil";

type SpecialHolidayFormProps = {
  editItem: SpecialHoliday;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (item: SpecialHoliday): void;
}
interface callbackCancel {
  (): void;
}

type SpecialHolidayInput = {
  id: string;
  name: string;
  start: Date;
  end: Date;
};

const convertInput = (item: SpecialHoliday): SpecialHolidayInput => {
  return {
    id: item.id,
    name: item.name,
    start: ToDate(item.start),
    end: ToDate(item.end),
  };
};

const reverseInput = (item: SpecialHolidayInput): SpecialHoliday => {
  return {
    id: item.id,
    name: item.name,
    start: ToDateString(item.start),
    end: ToDateString(item.end),
  };
};

const errorMessage = ({
  name,
  error,
  maxLength,
  min,
  max,
}: {
  name: string;
  error: FieldError | undefined;
  maxLength?: string;
  min?: string;
  max?: string;
}): string => {
  if (error?.type === "required") {
    return "入力が必要です。";
  }
  if (error?.type === "maxLength") {
    return maxLength + "文字以下にしてください。";
  }
  if (error?.type === "min") {
    return min + "以上にしてください。";
  }
  if (error?.type === "max") {
    return max + "以下にしてください。";
  }
  return "";
};

export default function SpecialHolidayForm(props: SpecialHolidayFormProps) {
  const {
    register,
    control,
    handleSubmit,
    setError,
    setValue,
    formState: { errors },
  } = useForm<SpecialHolidayInput>({
    defaultValues: convertInput(props.editItem),
  });
  const onSubmit: SubmitHandler<SpecialHolidayInput> = (data) => {
    const converted = reverseInput(data);
    if (!startDateIsLessThanEndDate(converted.start, converted.end)) {
      setError(`start`, { message: "開始日時が終了時刻よりも大きいです。" });
      setError(`end`, { message: "開始日時が終了時刻よりも大きいです。" });
      return;
    }
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
            {...register("name", { required: true, maxLength: 10 })}
            error={Boolean(errors.name)}
            helperText={errorMessage({
              name: "名称",
              error: errors.name,
              maxLength: "10",
            })}
          />
          <Controller
            name="start"
            control={control}
            render={({ field }) => {
              return (
                <MobileDatePicker
                  {...field}
                  label="開始日"
                  mask="____/__/__ __:__:__"
                  inputFormat="yyyy/MM/dd"
                  renderInput={(props) => (
                    <TextField
                      {...props}
                      error={Boolean(errors.start)}
                      helperText={errorMessage({
                        name: "開始日",
                        error: errors.start,
                      })}
                    />
                  )}
                  onChange={(newValue) => {
                    if (newValue != null) {
                      setValue("start", newValue);
                    }
                  }}
                />
              );
            }}
          />
          <Controller
            name="end"
            control={control}
            render={({ field }) => {
              return (
                <MobileDatePicker
                  {...field}
                  label="終了日"
                  mask="____/__/__ __:__:__"
                  inputFormat="yyyy/MM/dd"
                  renderInput={(props) => <TextField {...props} />}
                  onChange={(newValue) => {
                    if (newValue != null) {
                      setValue("end", newValue);
                    }
                  }}
                />
              );
            }}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            確定
          </Button>
          <Button
            color="error"
            variant="contained"
            size="large"
            onClick={onCancel}
          >
            キャンセル
          </Button>
        </Stack>
      </Container>
    </>
  );
}
