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
		`Our Dockerfile will have two section, first one where we build the binary and the second one which will be our final image. `,
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
func BenchmarkBaseIndex(t *testing.B)  {
	bi := NewBaseIndex()
	bi.Add( documents... )

	t.ResetTimer()

	for i := 0; i < t.N; i ++ {
		for _,word := range []string{`Нет подключения к Интернету`, `Choosing Go is a wise decision`} {
			bi.Find(word)
		}
	}
}

//
func BenchmarkCharsIndex(t *testing.B)  {

	ic := NewIndexChars()
	ic.Add(documents...)

	t.ResetTimer()

	for i := 0; i < t.N; i ++ {
		for _,word := range []string{`Нет подключения к Интернету`, `Choosing Go is a wise decision`} {
			ic.Find(word)
		}
	}
}

//
func TestBaseIndex(t *testing.T)  {
	bi := NewBaseIndex()
	bi.Add(
		`Нет подключения к Интернету`,
		`Create container from the image and expose it by mentioning a port`,
		`Please consider chucking me a couple of quid for my time and effort.`,
		`You can also perform actions on individual containers. This example prints the logs of a container given its ID.`,
	)
	tFindPositive(t, bi, `Нет`)
	tFindPositive(t, bi, `к`)
	tFindPositive(t, bi, `Create`)
	tFindPositive(t, bi, `individual`)
	tFindPositive(t, bi, `This`)
	tFindPositive(t, bi, `eXample`)
	tFindPositive(t, bi, `on`)
	tFindNegative(t, bi, `php`)
	tFindNegative(t, bi, `t`)
	tFindNegative(t, bi, `а`)
	tFindNegative(t, bi, `и`)
}

//
func TestCharsIndex(t *testing.T)  {
	bi := NewIndexChars()
	bi.Add(
		`Нет подключения к Интернету`,
		`Иванова нет на месте`,
		`Create container from the image and expose it by mentioning a port`,
		`Please consider chucking me a couple of quid for my time and effort.`,
		`You can also perform actions on individual containers. This example prints the logs of a container given its ID.`,
	)
	tFindPositive(t, bi, `Нет`)
	tFindPositive(t, bi, `к`)
	tFindPositive(t, bi, `Create`)
	tFindPositive(t, bi, `individual`)
	tFindPositive(t, bi, `This`)
	tFindPositive(t, bi, `eXample`)
	tFindPositive(t, bi, `on`)
	tFindNegative(t, bi, `php`)
	tFindNegative(t, bi, `t`)
	tFindNegative(t, bi, `а`)
	tFindNegative(t, bi, `и`)

	tFindNegative(t, bi, `и`)
	tFindNegative(t, bi, `Интернет`)
	tFindPositive(t, bi, `Интернет(у)`)
	tFindPositive(t, bi, `Иванов(а|ой|ым)`)
}


//
func TestBaseIndexMath(t *testing.T)  {
	bi := NewBaseIndex()
	bi.Add(
		`Нет подключения к Интернету`,
		`Иванова нет на месте`,
	)
	tFindNegative(t, bi, `и`)
	tFindNegative(t, bi, `Интернет`)
	tFindPositive(t, bi, `Интернет(у)`)
	tFindPositive(t, bi, `Иванов(а|ой|ым)`)
}

//
func tFindPositive(t *testing.T, i Index, word string) {
	if i.Find(word) == false {
		t.Fatalf(`Error not find word %s`, word)
	}
}

//
func tFindNegative(t *testing.T, i Index, word string) {
	if i.Find(word) == true {
		t.Fatalf(`Error wrong find word %s`, word)
	}
}