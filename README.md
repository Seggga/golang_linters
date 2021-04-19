# golang_linters
домашка к уроку 6 - Линтеры

1) deadcode нашел функцию, которую я создавал для тестирования на реальной файловой системе.
После перехода на testing/fstest она стала невостребованной.
Функцию из кода убрал.

cloremover\cloremover_stab.go:148:6: `createTestFiles` is unused (deadcode)
func createTestFiles() {
     ^

2) gosimple предложил упростить условие, менять не стал. Полагаю, так читабельнее.

cloremover\read_flags.go:26:5: S1002: should omit comparison to bool constant, can be simplified to `!*removeFlag` (gosimple)
        if *removeFlag == false && *confirmFlag == "off" {
           ^

3) unparam нашел неиспользуемый аргумент при вызове функции enumDirs - удалил.

cloremover\find_clones.go:79:15: `enumDirs` - `conf` is unused (unparam)
func enumDirs(conf *ConfigType, fileSystem fs.FS) ([]string, error) {
              ^

4) gosec рекомендует изменить права на доступ к каталогу для анализа: вместо 666 выставить 600. Исправил.

main.go:22:12: G302: Expect file permissions to be 0600 or less (gosec)
        f, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.SEEK_CUR|os.O_APPEND, 0666)
                  ^

5) errorlint рекомендует обертывать ошибки с помощью глагола %w. Исправил в четырех местах.

cloremover\remove_clones.go:81:64: non-wrapping format verb for fmt.Errorf. Use `%w` to format errors (errorlint)
                return 0, fmt.Errorf("There is an error entering data.\n%v", err)
                                                                             ^

