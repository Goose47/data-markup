export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>;

export type MarkupTypeField = {
  id?: number;
  name: string;
  label: string;
  assessment_type_id: number;
  group_id: number;
};

export type MarkupTypeRq = {
  name: string;
  fields: MarkupTypeField[];
};

export type MarkupTypeUpdateRq = MarkupTypeRq & {
  batch_id: number;
};

export type MarkupType = {
  id: number;
  batch_id: number | null;
  name: string;
  user_id: number | null;
  fields?: MarkupTypeField[] | null;
  markup_count: number;
  assessment_count: number;
  correct_assessment_count: number;
  created_at: string;
};

export type BatchRq = {
  name: string;
  overlaps: number;
  priority: number;
  markups: File;
  type_id: number;
};

export type BatchCardType = {
  id: number;
  name: string;
  overlaps: number;
  priority: number;
  created_at: string;
  is_active: boolean;
  type_id: number;
};

export type BatchMarkupType = {
  id: number;
  batch_id: number;
  status_id: number;
  data: string;
  correct_assessment_hash: string | null;
  assessments: any[];
  markup_count: number;
  processed_markup_count: number;
  assessment_count: number;
  correct_assessment_count: number;
};

export type AssessmentNext = {
  assessment_id?: number;
  /** assume it's a valid JSON */
  data: string;
  markup_type: MarkupType;
};

export type AssesmentUpdateField = {
  text: string | null;
  markup_type_field_id: number;
};

export type AssessmentUpdateRq = {
  fields: AssesmentUpdateField[];
};

export type AssessmentStoreRq = AssessmentUpdateRq & {
  markup_id: number;
};

export type MarkupTypeFull = RequiredFields<MarkupType, "fields">;

export type UserType = {
  id: number;
  email: string;
  created_at: string;
  roles: any[] | null;
};

export type AsseessmentType = {
  id: number;
  user: UserType;
  is_prior: boolean;
  created_at: string;
  updated_at: string;
  hash: string;
  markup_id?: number;
  fields: {
    id: number;
    assessment_id: number;
    markup_type_field_id: number;
    text: string | null;
    markup_type_field?: {
      id: number;
      name: string;
      label: string;
    };
  }[];
  markup_type: {
    fields: {
      id: number;
      name: string;
      label: string;
    }[];
  };
  is_correct?: boolean;
  is_editable?: boolean;
};

export type UserListType = {
  id: number;
  email: string;
  created_at: string;
  assessment_count: number;
  correct_assessment_count: number;
  assessment_count2: number;
  correct_assessment_count2: number;
};
