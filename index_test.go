package word_index

import (
	"testing"
)

var (
	documents = []string{
		`Sorry. We’re having trouble getting your pages back.`,
		`Наши рестораны это необычный формат: открытая кухня, блюда со всего света, фермерские продукты, уют и дружеская атмосфера.`,
		`If that address is correct, here are three other things you can try: Try again later.    Check your network connection.    If you are connected but behind a firewall, check that Firefox has permission to access the Web.`,
		`Still not able to restore your session? Sometimes a tab is causing the issue. View previous tabs, remove the checkmark from the tabs you don’t need to recover, and then restore.`,
		`We are having trouble restoring your last browsing session. Select Restore Session to try again.`,
		`Docker supports multi-stage builds, meaning docker will build in one container and you can copy the build artifacts to final image.`,
		`Нет подключения к Интернету`,
		`В качестве исходных данных будет использована статистика из Яндекс.Метрики. DataLens автоматически создаст дашборд на основе счетчика Метрики с подборкой графиков, а вы сможете отредактировать его по своему усмотрению.`,
		`Войти в личный кабинет ePayments`,
		`В математическом анализе и информатике кривая Мортона, Z-последовательность,Z-порядок, кривая Лебега, порядок Мортона или код Мортона — это функция, которая отображает многомерные данные в одномерные, сохраняя локальность точек данных. Функция была введена в 1966 Гаем Макдональдом Мортоном[1].`,
		`Our Dockerfile will have two section, first one where we build the binary and the second one which will be our final image. `,
		`Мы рассмотрим вашу заявку и ответим вам в ближайшее время. Если заявка будет подтверждена, доступ к сервису появится автоматически.`,
		`If you want to keep your array ordered, you have to shift all of the elements at the right of the deleting indexRegexp by one to the left. Hopefully, this can be done easily in Golang:`,
		`Данный шаг доступен для пользователей, у которых есть права на какой-либо счетчик Метрики. Если у вас нет прав на счетчик, то откройте готовый дашборд Metriсa и перейдите к шагу 2.`,
		`Удачное расположение с панорамой Казанского собора, вежливый и приветливый персонал, очень вкусная еда, быстрое обслуживание создали хорошее настроение и комфорт. `,
		`Ресторан быстрого обслуживания Marketplace (бывшие «Фрикадельки») — это новая, современная интерпретация демократичного ресторана с открытой кухней, линией раздачи «Free flow» и отделом кулинарии.`,
		`Marketplace – это демократичный ресторан с открытой кухней и живой атмосферой европейского рынка.`,
		`Георгию Карамзину и Татьяне Самсоновой предъявлено обвинение как посредникам. Изначально и Гоголев, и Карамзин были арестованы, но накануне Якутский городской суд перевел последнего под домашний арест. Самсонова же под домашним арестом изначально.`,
		`Гоголев — главный обвиняемый по этому делу. По данным следствия, девелопер передал ему взятку в виде прав на недвижимость в строящемся жилом доме на 1491 квадратный метр.`,
		`Попробуйте сделать следующее: Проверьте сетевые кабели, модем и маршрутизатор. Подключитесь к сети Wi-Fi ещё раз.`,
		`This is almost identical to the first Makefile we created for our consignment-service, however notice the service names and the ports have changed a little`,
		`«В этой связи суд отказал в изменении меры пресечения на заключение под стражу, считая выявленные нарушения недостаточными для немедленного изменения меры пресечения», — рассказали в пресс-службе.`,
		`You can also perform actions on individual containers. This example prints the logs of a container given its ID.`,
		`Note: Don’t run this on a production server. Also, if you are using swarm services, the containers stop, but Docker creates new ones to keep the service running in its configured state.`,
		`This first example shows how to run a container using the Docker API. On the command line, you would use the docker run command, but this is just as easy to do from your own apps too.`,
		`Each of these examples show how to perform a given Docker operation using the Go and Python SDKs and the HTTP API using curl.`,
		`Create container from the image and expose it by mentioning a port`,
		`Please consider chucking me a couple of quid for my time and effort.`,
		`Choosing Go is a wise decision that gives scalability and concurrency for your application and selecting a light weight image like alpine will make the push and pull of the image to registries faster, also small size base gives you minimal operating features to build functional container where you can add/install necessary dependencies in future`,
		`“Golang” the language created by Google that provides impeccable performance over application that demands concurrency has grabbed a separate spot in the tech community, few well known Inc’s that adopted the language include Facebook, Netflix , Dropbox etc.`,
	}
)

//
func BenchmarkIndexRegexp(t *testing.B) {
	bi := NewIndexRegexp()
	bPlainText(t, bi)
}

//
func BenchmarkIndexBin(t *testing.B) {
	ic := NewIndexBin()
	bPlainText(t, ic)
}

//
func BenchmarkIndexInterpolation(t *testing.B) {
	ic := NewIndexInterpolation()
	bPlainText(t, ic)
}

//
func bPlainText(t *testing.B, i Index) {
	i.Add(documents...)

	t.ResetTimer()

	for j := 0; j < t.N; j++ {
		for _, word := range []string{`from the image and expose it by`, `my time and effort`, `Нет подключения к Интернету`, `Choosing Go is a wise decision`} {
			i.Find(word)
		}
	}
}

//
func TestIndexRegexp(t *testing.T) {
	bi := NewIndexRegexp()
	tIndexPlainText(t, bi)
	tIndexMathText(t, bi)
}

//
func TestIndexBin(t *testing.T) {
	bi := NewIndexBin()
	tIndexPlainText(t, bi)
	tIndexMathText(t, bi)
}

//
func TestIndexBinAtDocument(t *testing.T) {
	bi := NewIndexBin()
	tAtDocument(t, bi)
}

//
func TestIndexInterpolationAtDocument(t *testing.T) {
	bi := NewIndexInterpolation()
	tAtDocument(t, bi)
}

//
func TestIndexRegexpAtDocument(t *testing.T) {
	bi := NewIndexRegexp()
	tAtDocument(t, bi)
}

//
func TestIndexInterpolation(t *testing.T) {
	bi := NewIndexInterpolation()
	tIndexPlainText(t, bi)
	tIndexMathText(t, bi)
}

//
func tAtDocument(t *testing.T, bi Index) {
	ss := []string{`first str`, ``, `third str`}
	bi.Add(ss...)

	for i, s := range ss {
		str, ok := bi.DocumentAt(i)
		if !ok {
			t.Fatalf(`Index %d not extists`, i)
		}
		if str != s {
			t.Fatalf(`String not equals: %s %s`, str, s)
		}
	}

	_, ok := bi.DocumentAt(len(ss) + 1)
	if ok {
		t.Fatalf(`Index %d out of range`, len(ss)+1)
	}
}

//
func tIndexMathText(t *testing.T, i Index) {
	i.Add(
		`Нет подключения к Интернету`,
		`Иванова нет на месте`,
		`Create container from the image and expose it by mentioning a port`,
	)
	tFindNegative(t, i, `и`)
	tFindNegative(t, i, `Интернет`)
	tFindPositive(t, i, `Интернет(у)`)
	tFindPositive(t, i, `Иванов(а|ой|ым)`)
	tFindNegative(t, i, `contai`)
	tFindPositive(t, i, `contai(ner)`)
	tFindPositive(t, i, `contai*`)
	tFindPositive(t, i, `contai(xxx|*)`)
	tFindPositive(t, i, `c*`)
	tFindPositive(t, i, `mentioning*`)
}

//
func tIndexPlainText(t *testing.T, i Index) {
	i.Add(
		`Нет подключения к Интернету`,
		`Create container from the image and expose it by mentioning a port`,
		`Please consider chucking me a couple of quid for my time and effort.`,
		`You can also perform actions on individual containers. This example prints the logs of a container given its ID.`,
	)
	tFindPositive(t, i, `Нет`)
	tFindPositive(t, i, `к`)
	tFindPositive(t, i, `Create`)
	tFindPositive(t, i, `individual`)
	tFindPositive(t, i, `This`)
	tFindPositive(t, i, `eXample`)
	tFindPositive(t, i, `on`)
	tFindNegative(t, i, `php`)
	tFindNegative(t, i, `t`)
	tFindNegative(t, i, `а`)
	tFindNegative(t, i, `и`)
}

//
func tFindPositive(t *testing.T, i Index, word string) {
	if i.Find(word) == -1 {
		t.Fatalf(`Error not find word %s`, word)
	}
}

//
func tFindNegative(t *testing.T, i Index, word string) {
	if i.Find(word) >= 0 {
		t.Fatalf(`Error wrong find word %s`, word)
	}
}
