import { Checkbox, Radio, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { MarkupTypeField } from "../../utils/types";
import "./MarkupTypeGroup.scss";
import { useMemo } from "react";
import { FieldValue } from "../MyMarkupType/MyMarkupType";

const _ = require("lodash");

const b = block("markup-type-group");

export const MarkupTypeGroup = ({
  fields,
  onUpdate,
  value,
}: {
  fields: MarkupTypeField[];
  onUpdate: (value: FieldValue[]) => void;
  value?: FieldValue[];
}) => {
  const uniqueId = useMemo(() => _.uniqueId(), []);

  if (fields.length === 0 || value === undefined) {
    return <></>;
  }

  if (fields[0].assessment_type_id === 1) {
    return (
      <div className={b("radio")}>
        <label className={b("label")}>{fields[0].label}</label>
        {fields.map((field) => (
          <Radio
            id={`${uniqueId}_${String(field.group_id)}_${String(field.id)}`}
            value={String(field.id)}
            content={field.name}
            onUpdate={(checked) => {
              let value: string = '';
              if (checked) {
                value = String(field.id);
              }
              onUpdate([ {
                value: value,
                assessment_type_id: field.assessment_type_id,
              } ]);
            }}
            checked={value[0].value === String(field.id)}
          />
        ))}
      </div>
    );
  } else if (fields[0].assessment_type_id === 2) {
    return (
      <div className={b("checkbox")}>
        <label className={b("label")}>{fields[0].label}</label>
        {fields.map((field) => (
          <Checkbox
            id={`${uniqueId}_${String(field.group_id)}_${String(field.id)}`}
            value={String(field.id)}
            content={field.name}
            onUpdate={(checked) => {
              let valueCopy: FieldValue[] = _.cloneDeep(value);
              if (checked) {
                if (!valueCopy.some(({value}) => value === String(field.id))) {
                  valueCopy.push({value: String(field.id), assessment_type_id: field.assessment_type_id});
                }
              } else {
                valueCopy.splice(valueCopy.findIndex(({ value }) => value === String(field.id)), 1);
              }
              onUpdate(valueCopy);
            }}
            checked={value.some(({value}) => value === String(field.id))}
          />
        ))}
      </div>
    );
  } else if (fields[0].assessment_type_id === 3) {
    return (
      <div className={b("select")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <Select
          placeholder={fields[0].label}
          width={"max"}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          value={value.map(({ value }) => value)}
          onUpdate={(value) => onUpdate(value.map((v) => ({ value: v, assessment_type_id: fields[0].assessment_type_id })))}
        >
          {fields.map((field) => (
            <Select.Option value={String(field.id)}>{field.name}</Select.Option>
          ))}
        </Select>
      </div>
    );
  } else if (fields[0].assessment_type_id === 4) {
    return (
      <div className={b("select")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <Select
          multiple={true}
          placeholder={fields[0].label}
          width={"max"}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          value={value.map(({ value }) => value)}
          onUpdate={(value) => onUpdate(value.map((v) => ({ value: v, assessment_type_id: fields[0].assessment_type_id })))}
        >
          {fields.map((field) => (
            <Select.Option value={String(field.id)}>{field.name}</Select.Option>
          ))}
        </Select>
      </div>
    );
  } else if (fields[0].assessment_type_id === 5) {
    return (
      <div className={b("text")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <TextInput
          value={value[0].value}
          onUpdate={(v) => {
            onUpdate([{value: v, assessment_type_id: fields[0].assessment_type_id}]);
          }}
          placeholder={fields[0].label}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
        />
      </div>
    );
  }
  return <></>;
};
