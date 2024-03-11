export type SectionType =
  | 'TEST'
  | 'EXPECT'
  | 'EXPECTF'
  | 'EXPECTREGEX'
  | 'EXPECTHEADERS'
  | 'POST'
  | 'POST_RAW'
  | 'GZIP_POST'
  | 'DEFLATE_POST'
  | 'PUT'
  | 'GET'
  | 'COOKIE'
  | 'ARGS'
  | 'FILE'
  | 'CAPTURE_STDIO'
  | 'STDIN'
  | 'CGI'
  | 'INI'
  | 'ENV'
  | 'EXTENSIONS'
  | 'SKIPIF'
  | 'XFAIL'
  | 'CLEAN'
  | 'CREDITS'
  | 'DESCRIPTION'
  | 'CONFLICTS'
  | 'WHITESPACE_SENSITIVE'

export type SectionMap = { [K in SectionType]?: string }
export type Section = { type: SectionType; text: string }

// section priorities
const sectionTypePriority: { [K in SectionType]: number } = {
  TEST: 1,
  EXPECT: 2,
  EXPECTF: 3,
  EXPECTREGEX: 4,
  EXPECTHEADERS: 5,
  POST: 6,
  POST_RAW: 7,
  GZIP_POST: 8,
  DEFLATE_POST: 9,
  PUT: 10,
  GET: 11,
  COOKIE: 12,
  ARGS: 13,
  FILE: 14,
  CAPTURE_STDIO: 15,
  STDIN: 16,
  CGI: 17,
  INI: 18,
  ENV: 19,
  EXTENSIONS: 20,
  SKIPIF: 21,
  XFAIL: 22,
  CLEAN: 23,
  CREDITS: 24,
  DESCRIPTION: 25,
  CONFLICTS: 26,
  WHITESPACE_SENSITIVE: 27
}

export function compareSectionType(t1: SectionType, t2: SectionType): number {
  const p1 = sectionTypePriority[t1] || 999
  const p2 = sectionTypePriority[t2] || 999
  return p1 - p2
}

// convert sectionMap to sections
export function sectionMapToList(m: SectionMap): Section[] {
  const sections: Section[] = Object.keys(m).map((key) => {
    const type = key as SectionType
    const text = m[type] + ''
    return { type, text }
  })
  return sections.sort((sec1, sec2) => compareSectionType(sec1.type, sec2.type))
}

// convert sections to sectionMap
export function sectionListToMap(list: Section[]): SectionMap {
  const m: SectionMap = {}
  list.forEach((sec) => {
    m[sec.type] = sec.text
  })
  return m
}
